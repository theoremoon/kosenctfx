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
from minio import Minio
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
    self._manager = API(
      self._conf["manager"]["url"],
      self._conf["manager"]["token"],
    )
    self._do = API(
      self._conf["digitalocean"]["url"],
      self._conf["digitalocean"]["token"],
    )

    # uploaderの設定。Minioを使ってS3 compatibleなendpointを利用できる
    bucket = self._conf["bucket"]
    self._minio = Minio(bucket["endpoint"], access_key=bucket["access_key"], secret_key=bucket["secret_key"], secure=False if "insecure"  in bucket else True)

  def manager_init(self):
    r = self._manager.post("/init", {
      "server_url": self._conf["scoreserver"]["url"],
      "server_token": self._conf["scoreserver"]["token"],
    })
    print("[+] challenge manager initialized")

  def bucket_init(self):
    # バケットがなければ作り、ついでにポリシーを制定する
    # よく考えるとこのツールでbucketを作ってるのはおかしいんだよな……
    # おかしくはなくて、bucket名の指定とAPI keyを渡しているのは正しい（？）
    # どうせwrite keyは必要だからね……
    # どちらかといえばDigital Ocean側の鍵を持っているのがおかしい
    # Read可能 / List不可能 / Write不可能
    bucket = self._conf["bucket"]
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
      print("[-] {}".format(r.json()))

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
          "IsRunning": rc["is_running"],
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

  def close_firewall(self, ids):
    raise Exception("unimplemented")

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

  def open_firewall(self, ids):
    raise Exception("unimplemented")

  def start(self, ids):
    """
    - 問題を起動する。idsは数値またはリスト
    """
    if isinstance(ids, int):
      ids = [ids]
    else:
      ids = ids

    for id in ids:
      r = self._manager.post("/start", {
        "id": id
      })
      if 200 <= r.status_code < 400:
        print("[+] started: {}".format(id))
      else:
        print("[-] failed to start: {}".format(id))
        print(r.text)

  def stop(self, ids):
    """
    - 問題を停止する。idsは数値またはリスト
    """
    if isinstance(ids, int):
      ids = [ids]
    else:
      ids = ids

    for id in ids:
      r = self._manager.post("/stop", {
        "id": id
      })
      if 200 <= r.status_code < 400:
        print("[+] stopped: {}".format(id))
      else:
        print("[-] failed to stop: {}".format(id))
        print(r.text)

  def _upload_file(self, fileobj, size, name):
      self._minio.put_object(self._conf["bucket"]["name"], name, fileobj, size)

      # uploadしたファイルのdownloadableなリンク
      url = "{schema}://{host}/{name}/{key}".format(
          schema="http" if "insecure" in self._conf["bucket"] else "https",
          host=self._conf["bucket"]["endpoint"],
          name=self._conf["bucket"]["name"],
          key=name,
      )
      return url


  def _compress_distfiles(self, dir):
      # 圧縮パート
      tarbytes = BytesIO()
      with tarfile.open(fileobj=tarbytes, mode="w:gz") as tar:
        for f in dir.iterdir():
            tar.add(f, arcname=f.as_posix()[len(dir.as_posix()):], recursive=True)
      return tarbytes.getvalue()

  def _build_push_image(self, compose_path, image_name):
      """
      既存のdocker-compose.ymlを読み込んで、imageセクションを追加してbuild / pushし、新しくなったdocker-compose.ymlの内容を返す
      """
      dir = compose_path.parent

      # 既存のdocker-compose.ymlを読み込んで imageを追加する
      with open(compose_path) as f:
        compose = yaml.safe_load(f)
      for service in compose["services"].keys():
        if "image" in service:
          print("[+] unsupported!!!")
          quit()

        # imageタグをセットする
        compose["services"][service]["image"] = "{}/{}/{}_{}:latest".format(self._conf["registry"]["server"], self._conf["registry"]["name"], image_name, service)

        # buildは削除する
        if "build" in service:
          del compose["services"][service]["build"]

      # 新しいcompose.ymlを作る
      new_compose_path = dir / "docker-compose_{}.yml".format(randname())
      new_compose = yaml.dump(compose, default_flow_style=False)
      with open(new_compose_path, "w") as f:
        f.write(new_compose)

      # 新しく作ったdocker-compose.ymlを使ってイメージをビルド、pushする
      subprocess.run(["docker-compose", "-f", new_compose_path.name, "build"], cwd=dir, check=True)
      subprocess.run(["docker", "login", self._conf["registry"]["server"], "-u", self._conf["registry"]["username"], "-p", self._conf["registry"]["password"]], check=True)
      subprocess.run(["docker-compose", "-f", new_compose_path.name, "push"], cwd=dir, check=True)

      # 新しく作ったdocker-compose.ymlはもういらないので消す
      os.remove(new_compose_path)

      # 新しいdocker-composeの内容だけを返す
      return new_compose

  def _register_to_manager(self, compose_file, challengeid, taskinfo):
      new_compose = self._build_push_image(compose_file, md5(taskinfo["name"].encode()).hexdigest())
      solve_compose_path = compose_file.parent / "solution" / "docker-compose.yml"
      if solve_compose_path.is_file():
        solve_compose = do_compose_iikanji(solve_compose_path, md5((taskinfo["name"] + "_solution").encode()).hexdigest())
      else:
        solve_compose = ""

      # challenge managerに登録する
      r = self._manager.post("/register", {
        "id": challengeid,
        "compose": new_compose,
        "solve_compose": solve_compose,
        "flag": taskinfo["flag"],
        "host": taskinfo["host"],
        "port": taskinfo["port"],
      })
      if 200 <= r.status_code < 400:
          print("[+] registered to manager: {}".format(taskinfo["name"]))
      else:
          print("[-] {}".format(r.text))

  def register(self, no_manager=False):
    """
    - 問題を全てスコアサーバに登録する
    - 配布ファイルをサーバにアップロードする
    - Docker registryに問題イメージと問題solveイメージをアップロードする
    """

    ## 問題をスコアサーバに登録する
    for chal in iterate_challenges(self._basedir):
      with open(chal / "task.json", "r") as f:
        taskinfo = json.load(f)

      if "port" in taskinfo:
        taskinfo["port"] = str(taskinfo["port"])
      # hostとかportとかをdescriptionに埋め込んでいる場合
      taskinfo["description"] = Template(taskinfo["description"]).substitute(taskinfo)
      taskinfo["attachments"] = []

      # distfiles以下のファイルを圧縮してupload
      distdir = chal / "distfiles"
      if distdir.is_dir():
        tarbytes = self._compress_distfiles(distdir)
        object_name = "{}_{}.tar.gz".format(taskinfo["name"], hashlib.md5(tarbytes).hexdigest())
        url = self._upload_file(BytesIO(tarbytes), len(tarbytes), object_name)
        taskinfo["attachments"].append({
          "name": object_name,
          "url": url,
        })

      # rawdistfiles以下のファイルはそのままupload
      rawdistdir = chal / "rawdistfiles"
      if rawdistdir.is_dir():
        for f in rawdistdir.iterdir():
          if f.is_file():
            with open(f) as fileobj:
              url = self._upload_file(fileobj, f.stat().st_size, f.name)
            taskinfo["attachments"].append({
              "name": f.name,
              "url": url,
            })

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
      challengeinfo = r.json()

      compose_file = chal / "docker-compose.yml"
      if no_manager == False and compose_file.is_file():
        self._register_to_manager(compose_file, challengeinfo["model"]["id"], taskinfo)

def main():
  fire.Fire(CommandClass)

if __name__ == '__main__':
  main()
