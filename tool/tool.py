#!/usr/bin/env python3
import os
import json
import glob
import fire
from pathlib import Path
import requests
import pickle
import configparser
import tarfile
from io import BytesIO
from minio import Minio
import hashlib
from string import Template
import tableprint

BASEDIR_DEFAULT = os.getcwd()
CONFFILE_DEFAULT = Path(os.getcwd()) / "config.ini"

def login(server, username, password):
  """
  スコアサーバにログインしてSessionを取得し、 cookie.pickleとして保存する
  """
  r = requests.post(server + "/login", data=json.dumps({
    "username": username,
    "password": password,
  }), headers={"Content-Type": "application/json"})
  r.raise_for_status()

  return r.cookies


def iterate_challenges(basedir):
  """
  challengesディレクトリを探索して問題一覧を表示する
  """
  basepath = Path(basedir)
  for taskdir in basepath.glob("**/task.json"):
    yield taskdir.parent


class API():
  def __init__(self, url, username, password):
    self.url = url
    self.username = username
    self.password = password
    self.cookies = login(url, username, password)

  def post(self, endpoint, data):
    while True:
      r = requests.post(self.url + endpoint, data=json.dumps(data), headers={"Content-Type": "application/json"}, cookies=self.cookies)
      if r.status_code == 401:
        self.cookies = login(self.url, self.username, self.password)
        continue

      return r

  def get(self, endpoint):
    while True:
      r = requests.get(self.url + endpoint, cookies=self.cookies)
      if r.status_code == 401:
        self.cookies = login(self.url, self.username, self.password)
        continue

      return r


class CommandClass():
  def __init__(self, basedir=BASEDIR_DEFAULT, configfile=CONFFILE_DEFAULT):
    self._basedir = basedir
    self._conf = configparser.ConfigParser()
    self._conf.read(configfile)
    self._api = API(
      self._conf["scoreserver"]["url"],
      self._conf["scoreserver"]["username"],
      self._conf["scoreserver"]["password"],
    )

    # uploaderの設定。Minioを使ってS3 compatibleなendpointを利用できる
    bucket = self._conf["bucket"]
    self._minio = Minio(bucket["endpoint"], access_key=bucket["access_key"], secret_key=bucket["secret_key"], secure=False if "insecure"  in bucket else True)

    # バケットがなければ作り、ついでにポリシーを制定する
    # Read可能 / List不可能 / Write不可能
    if not self._minio.bucket_exists(bucket["name"]):
      self._minio.make_bucket(bucket["name"], bucket["region"])
      self._minio.set_bucket_policy(bucket["name"], json.dumps({
        "Version":"2012-10-17",
        "Statement":[
          {
            "Sid":"AddPerm",
            "Effect":"Allow",
            "Principal": "*",
            "Action":["s3:GetObject"],
            "Resource":["arn:aws:s3:::{}/*".format(bucket["name"])]
          }
        ]
      }))

  def list(self):
    """
    問題一覧とそのステータスを表示する
    scoreserverのデータとローカルのデータを突き合わせる
    """

    r = self._api.get("/admin/list-challenges")
    remote = r.json()

    local = []
    for chal in iterate_challenges(self._basedir):
      with open(chal / "task.json", "r") as f:
        taskinfo = json.load(f)
      local.append(taskinfo)

    remote = {c["name"]:c for c in remote}
    local = {c["name"]:c for c in local}
    chals = []
    for name in remote.keys() | local.keys():
      lc = local.get(name)
      rc = remote.get(name)
      if rc:
        chals.append([
          rc["name"],
          "0x{:08x}".format(rc["id"]),
          rc["score"],
          str(rc["is_running"]), # FIXME
          str(rc["is_open"]),
          len(rc["solved_by"])
        ])
      elif lc:
        chals.append([lc["name"],"-","-","-","-","-"])

    tableprint.table(chals, [
      'Challenge',
      'ID',
      'Score',
      'IsRunning',
      'IsOpen',
      'Solved',
    ])


  def open(self, name):
    """
    - 問題を公開する。nameは文字列またはリスト
    """

    for chal in iterate_challenges(self._basedir):
      with open(chal / "task.json", "r") as f:
        taskinfo = json.load(f)
      if taskinfo["name"] != name or taskinfo["name"] not in name:
        continue

      if "host" in taskinfo:
        #TODO: Firewallを開ける
        pass

      r = self._api.post("/admin/open-challenge", {
        "name": taskinfo["name"],
      })
      print("[ ] {}: {}".format(taskinfo["name"], r.text))

  def start(self, name):
    """
    - 問題を起動する。nameは文字列またはリスト
    """
    pass


  def register(self):
    """
    - 問題を全てスコアサーバに登録する
    - 配布ファイルをサーバにアップロードする
    - Docker registryに問題イメージと問題solveイメージをアップロードする
    """

    ## 問題をスコアサーバに登録する
    for chal in iterate_challenges(self._basedir):
      with open(chal / "task.json", "r") as f:
        taskinfo = json.load(f)

      # hostとかportとかをdescriptionに埋め込んでいる場合
      taskinfo["description"] = Template(taskinfo["description"]).substitute(taskinfo)
      taskinfo["attachments"] = []

      # distfiles以下のファイルを圧縮してupload
      distdir = chal / "distfiles"
      if distdir.is_dir():
        # 圧縮パート
        tarbytes = BytesIO()
        with tarfile.open(fileobj=tarbytes, mode="w:gz") as tar:
          for f in distdir.iterdir():
              tar.add(f, arcname=f.as_posix()[len(distdir.as_posix()):], recursive=True)

        # uploadパート
        buf = tarbytes.getvalue()
        object_name = "{}_{}.tar.gz".format(taskinfo["name"], hashlib.md5(buf).hexdigest())
        tarbytes.seek(0)
        self._minio.put_object(self._conf["bucket"]["name"], object_name, tarbytes, len(buf))

        # uploadしたファイルのdownloadableなリンク
        url = "{schema}://{host}/{name}/{key}".format(
            schema="http" if "insecure" in self._conf["bucket"] else "https",
            host=self._conf["bucket"]["endpoint"],
            name=self._conf["bucket"]["name"],
            key=object_name,
        )
        taskinfo["attachments"].append({
          "name": object_name,
          "url": url,
        })

      # rawdistfiles以下のファイルはそのままupload
      rawdistdir = chal / "rawdistfiles"
      if rawdistdir.is_dir():
        for f in rawdistdir.iterdir():
          if f.is_file():
            self._minio.fput_object(self._conf["bucket"]["name"], f.name, f.as_posix())
            url = "{schema}://{host}/{name}/{key}".format(
                schema="http" if "insecure" in self._conf["bucket"] else "https",
                host=self._conf["bucket"]["endpoint"],
                name=self._conf["bucket"]["name"],
                key=f.name,
            )
            taskinfo["attachments"].append({
              "name": f.name,
              "url": url,
            })

      compose = chal / "docker-compose.yml"
      if compose.is_file():
        # TODO: yamlをいい感じにしてbuildしてpush
        pass

      r = self._api.post("/admin/new-challenge", data={
        "name": taskinfo["name"],
        "flag": taskinfo["flag"],
        "description": taskinfo["description"],
        "author": taskinfo["author"],
        "is_survey": taskinfo["is_survey"],
        "tags": taskinfo["tags"],
        "attachments": taskinfo["attachments"],
      })
      if 200 <= r.status_code < 400:
          print("[+] registered: {}".format(taskinfo["name"]))
      else:
          print("[-] {}".format(r.text))


def main():
  fire.Fire(CommandClass)

if __name__ == '__main__':
  main()
