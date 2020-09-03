#!/usr/bin/env python3
import os
import sys
import json
import glob
import fire
from pathlib import Path
import requests
import pickle
import configparser
import tarfile
from io import BytesIO
import boto3
from botocore.exceptions import ClientError
import hashlib
from string import Template
import yaml
import random
import string
from base64 import b64decode
from hashlib import md5
import subprocess
from datetime import datetime

BASEDIR_DEFAULT = os.getcwd()
CONFFILE_DEFAULT = Path(os.getcwd()) / "config.ini"

def randname():
  return "".join(random.choices(string.ascii_letters, k=8))


def iterate_challenges(basedir):
  """
  challengesディレクトリを探索して問題一覧を表示する
  """
  basepath = Path(basedir)
  for taskdir in basepath.glob("**/task.json"):
    yield taskdir.parent

class API():
  def __init__(self, url, token):
    self.url = url
    self.token = token

  def post(self, endpoint, data):
    r = requests.post(self.url + endpoint, data=json.dumps(data), headers={"Content-Type": "application/json", "Authorization": "Bearer {}".format(self.token)})
    return r

  def get(self, endpoint):
    r = requests.get(self.url + endpoint, headers={"Authorization": "Bearer {}".format(self.token)})
    return r


class CommandClass():
  def __init__(self, basedir=BASEDIR_DEFAULT, configfile=CONFFILE_DEFAULT):
    self._basedir = basedir
    self._conf = configparser.ConfigParser()
    self._conf.read(configfile)
    self._api = API(
      self._conf["scoreserver"]["url"],
      self._conf["scoreserver"]["token"],
    )

    # uploaderの設定。boto3を使ってS3 compatibleなendpointを利用できる
    bucket = self._conf["bucket"]
    session = boto3.session.Session()
    self._s3 = session.client("s3", region_name=bucket["region"], endpoint_url=("http://" if "insecure" in bucket else "https://") + bucket["endpoint"], aws_access_key_id=bucket["access_key"], aws_secret_access_key=bucket["secret_key"])

  def bucket_init(self):
    # バケットがなければ作り、ついでにポリシーを制定する
    # よく考えるとこのツールでbucketを作ってるのはおかしいんだよな……
    # おかしくはなくて、bucket名の指定とAPI keyを渡しているのは正しい（？）
    # どうせwrite keyは必要だからね……
    # どちらかといえばDigital Ocean側の鍵を持っているのがおかしい
    # Read可能 / List不可能 / Write不可能
    bucket = self._conf["bucket"]
    try:
      self._s3.head_bucket(Bucket=bucket["name"])
    except ClientError:
      self._s3.create_bucket(Bucket=bucket["name"])
      self._s3.put_bucket_policy(Bucket=bucket["name"], Policy=json.dumps({
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
    print("[+] bucket initialized")

  def apply_config(self, path):
    """
    CTFの設定を反映する
    """
    with open(path) as f:
      conf = yaml.safe_load(f)

    r = self._api.post("/admin/ctf-config", {
      "name": conf["ctf"]["name"],
      "start_at": int(datetime.strptime(conf["ctf"]["start"], "%Y-%m-%d %H:%M:%S %z").timestamp()),
      "end_at": int(datetime.strptime(conf["ctf"]["end"], "%Y-%m-%d %H:%M:%S %z").timestamp()),
      "register_open": conf["ctf"]["register_open"],
      "ctf_open": conf["ctf"]["ctf_open"],
      "lock_count": conf["ctf"]["lock_count_times"],
      "lock_duration": conf["ctf"]["lock_count_second"],
      "lock_second": conf["ctf"]["lock_second"],
      "score_expr": conf["ctf"]["score_expr"]
    })
    if 200 <= r.status_code < 400:
      print("[+] configured")
    else:
      print("[-] {}".format(r.text))

  def _local_challenges(self):
    chals = {}
    for chal in iterate_challenges(self._basedir):
      with open(chal / "task.json", "r") as f:
        taskinfo = json.load(f)
      chals[taskinfo["name"]] = taskinfo
    return chals

  def score_emulate(self, maxCount):
    """
    問題のsolve数と点数の関係をemulateする
    """
    r = self._api.get("/admin/score-emulate?maxCount={}".format(maxCount))
    if 200 <= r.status_code < 400:
      print("[+] configured")
    else:
      print("[-] {}".format(r.json()))
    print(r.json())

  def list(self):
    """
    問題一覧とそのステータスを表示する
    scoreserverのデータとローカルのデータを突き合わせる
    """

    r = self._api.get("/admin/list-challenges")
    remote = r.json()
    remote = {c["name"]:c for c in remote}

    local = self._local_challenges()
    chals = []
    for name in remote.keys() | local.keys():
      lc = local.get(name)
      rc = remote.get(name)
      if rc:
        chals.append({
          "Challenge": rc["name"],
          "ID": rc["model"]["id"],
          "Flag": rc["flag"],
          # "IsRunning": rc["is_running"],
          "IsOpen": rc["is_open"],
        })
      elif lc:
        chals.append({
          "Challenge": lc["name"],
          "Flag": lc["flag"],
        })

    print(json.dumps(chals, ensure_ascii=False, indent=4))



  def close(self, ids):
    """
    - 問題の公開をとりやめる。idsは数値またはリスト
    """
    if isinstance(ids, int):
      ids = [ids]
    else:
      ids = ids

    r = self._api.get("/admin/list-challenges")
    remote = r.json()
    for id in ids:
      for c in remote:
        if c["model"]["id"] != id:
          continue

        name = c["name"]

        # close する
        r = self._api.post("/admin/close-challenge", {
          "name": name,
        })
        if 200 <= r.status_code < 400:
          print("[+] closed: {}".format(name))
        else:
          print("[-] {}".format(r.text))
        break
      else:
        print("[-] no such challenge: {}".format(id))

  def open(self, ids):
    """
    - 問題を公開する。idsは数値またはリスト
    """
    if isinstance(ids, int):
      ids = [ids]
    else:
      ids = ids

    r = self._api.get("/admin/list-challenges")
    remote = r.json()
    for id in ids:
      for c in remote:
        if c["model"]["id"] != id:
          continue

        name = c["name"]
        # open する
        r = self._api.post("/admin/open-challenge", {
          "name": name,
        })
        if 200 <= r.status_code < 400:
          print("[+] opened: {}".format(name))
        else:
          print("[-] {}".format(r.text))
        break
      else:
        print("[-] no such challenge: {}".format(id))

  def _upload_file(self, fileobj, size, name):
      self._s3.put_object(Bucket=self._conf["bucket"]["name"], Key=name, Body=fileobj)

      # uploadしたファイルのdownloadableなリンク
      url = "{schema}://{host}/{name}/{key}".format(
          schema="http" if "insecure" in self._conf["bucket"] else "https",
          host=self._conf["bucket"]["endpoint"],
          name=self._conf["bucket"]["name"],
          key=name,
      )
      print("[+] uploaded: {}".format(url))
      return url

  def _compress_distfiles(self, dir):
      # 圧縮パート
      tarbytes = BytesIO()
      with tarfile.open(fileobj=tarbytes, mode="w:gz") as tar:
        for f in dir.iterdir():
            tar.add(f, arcname=f.as_posix()[len(dir.as_posix()):], recursive=True)
      return tarbytes.getvalue()


  def _task_upload(self, taskinfo, challengedir):
    attachments = []

    # distfiles以下のファイルを圧縮してupload
    distdir = challengedir / "distfiles"
    if distdir.is_dir():
      tarbytes = self._compress_distfiles(distdir)
      object_name = "{}_{}.tar.gz".format(taskinfo["name"], hashlib.md5(tarbytes).hexdigest())
      url = self._upload_file(BytesIO(tarbytes), len(tarbytes), object_name)
      attachments.append({
        "name": object_name,
        "url": url,
      })

    # rawdistfiles以下のファイルはそのままupload
    rawdistdir = challengedir / "rawdistfiles"
    if rawdistdir.is_dir():
      for f in rawdistdir.iterdir():
        if f.is_file():
          with open(f) as fileobj:
            url = self._upload_file(fileobj, f.stat().st_size, f.name)
          attachments.append({
            "name": f.name,
            "url": url,
          })
    return attachments


  def register(self, challenges=[], host=None):
    """
    - 問題を全てスコアサーバに登録する
    - 配布ファイルをサーバにアップロードする
    - Docker registryに問題イメージと問題solveイメージをアップロードする
    """

    ## 問題をスコアサーバに登録する
    for chal in iterate_challenges(self._basedir):
      with open(chal / "task.json", "r") as f:
        taskinfo = json.load(f)

      if len(challenges) > 0 and taskinfo["name"] not in challenges:
        continue

      if host and "host" in taskinfo:
        taskinfo["host"] = host

      if "port" in taskinfo:
        taskinfo["port"] = str(taskinfo["port"])
      # hostとかportとかをdescriptionに埋め込んでいる場合
      taskinfo["description"] = Template(taskinfo["description"]).substitute(taskinfo)
      taskinfo["attachments"] = self._task_upload(taskinfo, chal)

      # scoreserverに送る
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
          print("[+] registered to scoreserver: {}".format(taskinfo["name"]))
      else:
          print("[-] {}".format(r.text))

def main():
  fire.Fire(CommandClass)

if __name__ == '__main__':
  main()
