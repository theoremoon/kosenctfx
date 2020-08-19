from pathlib import Path
import string
import random
import os
import subprocess
import re
import json
import sys
import requests
import uuid
from base64 import b64decode

CONFIG_FILE = Path(__file__).parent / "config.json"
LOCAL_DIR = Path(__file__).parent / "boxes"
VAGRANTFILE = Path(__file__).parent / "Vagrantfile"

def get_config():
    with open(CONFIG_FILE, "r") as f:
        j = json.load(f)
    return j

def set_config(j):
    with open(CONFIG_FILE, "w") as f:
        json.dump(j, f)

class DigitalOcean():
    def __init__(self, config):
        self.config = config
        self.token = self.config["digitalocean"]["token"]

    def create_instance(self):
        name = gen_name()
        #TODO


class Local():
    def __init__(self, config):
        self.config = config

    def _random_ip(self):
        while True:
            ip = "192.168.33.{}".format(random.randint(1, 254))
            if ip not in self.config["local"]["used_ips"]:
                return ip

    def create_instance(self):
        # name / ip を決める
        name = gen_name()
        boxdir = LOCAL_DIR / name
        os.makedirs(boxdir)
        ip = self._random_ip()

        # box用のディレクトリを作ってup
        with open(VAGRANTFILE, "r") as f:
            vagrantfile = string.Template(f.read()).substitute({
                "ip": ip,
                "hostname": name,
            })
        with open(boxdir / "Vagrantfile", "w") as f:
            f.write(vagrantfile)
        subprocess.run(["vagrant", "up"], cwd=boxdir, check=True)

        # ssh-configを読んでくる
        r = subprocess.run(["vagrant", "ssh-config"], cwd=boxdir, capture_output=True, check=True)
        ssh_config = r.stdout.decode()

        with open(CONFIG_FILE, "w") as f:
            self.config["local"]["used_ips"].append(ip)
            json.dump(self.config, f)

        return (name, ip, ssh_config)

def gen_name():
    return "".join(random.choices(string.ascii_letters, k=16))

def make_instance():
    with open(CONFIG_FILE, "r") as f:
        config = json.load(f)

    if config["mode"] == "local":
        return Local(config)
    elif config["mode"] == "digitalocean":
        pass

def create_leader(box):
    # boxを作る
    name, ip, ssh_config = box.create_instance()
    config = get_config()

    # ansibleでdockerをinstallする
    with open("ssh_config", "w") as f:
        f.write(ssh_config)
    subprocess.run(["ansible-playbook", "-i", "hosts", "-e", "ip={}".format(ip), "ansible/docker_leader.yml"], check=True)

    # swarm init
    r = subprocess.run(["ssh", "-F", "ssh_config", "default", "docker swarm init --advertise-addr {ip} --listen-addr {ip}".format(ip=ip)], capture_output=True, check=True)
    join_command = re.findall(r"docker swarm join --token .+", r.stdout.decode())[0]
    with open("join_command", "w") as f:
        f.write(join_command)

    # leaderのIPを保存しておく
    config = get_config()
    config["leader"] = ip
    set_config(config)

    # ssh_config削除
    os.remove("ssh_config")


def create_client(box):
    # boxを作る
    name, ip, ssh_config = box.create_instance()
    config = get_config()

    # ansibleでdockerをinstallする
    with open("ssh_config", "w") as f:
        f.write(ssh_config)
    subprocess.run(["ansible-playbook", "-i", "hosts", "ansible/docker.yml"], check=True)

    # swarm join
    with open("join_command") as f:
        join_command = f.read()
    r = subprocess.run(["ssh", "-F", "ssh_config", "default", join_command], check=True)


    # ssh_config削除
    os.remove("ssh_config")

def create_manager(box):
    # boxを作る
    name, ip, ssh_config = box.create_instance()
    config = get_config()

    with open("ssh_config", "w") as f:
        f.write(ssh_config)

    # install docker and login to private registry
    subprocess.run(["ansible-playbook", "-i", "hosts", "ansible/docker.yml"], check=True)
    subprocess.run(["ssh", "-F", "ssh_config", "default", "sudo su -c 'docker login -u {username} -p {password} {server}'".format(**config["registry"])], check=True)

    # install challengemanager_tokenn
    challengemanager_token = str(uuid.uuid4())
    subprocess.run(["ansible-playbook", "-i", "hosts",
                    "-e", "docker_ip={}".format(config["leader"]),
                    "-e", "token={}".format(challengemanager_token),
                    "ansible/challengemanager.yml"], check=True)

    with open("manager.ini", "w") as f:
        f.write("""
[manager]
url=http://{}:5000
token={}
""".format(ip, challengemanager_token))

    # ssh_config削除
    os.remove("ssh_config")


def init(mode):
    with open("registry.json", "r") as f:
        registry = json.load(f)
        server, auth = list(registry["auths"].items())[0]
        username, password = b64decode(auth["auth"]).decode().split(":")
        registry = {
            "server": server,
            "username": username,
            "password": password,
        }
    if mode == "local":
        set_config({
            "local": {
                "used_ips": []
            },
            "mode": "local",
            "registry": registry,
        })
    elif mode == "digitalocean":
        set_config({
            "digitalocean": {
                "token": sys.argv[2]
            },
            "mode": "digitalocean",
            "registry": registry,
        })
    else:
        print("[-] init (local|digitalocean)")
        quit()


def main():
    if sys.argv[1] == "init":
        init(sys.argv[2])
    elif sys.argv[1] == "leader":
        box = make_instance()
        create_leader(box)
    elif sys.argv[1] == "client":
        box = make_instance()
        create_client(box)
    elif sys.argv[1] == "manager":
        box = make_instance()
        create_manager(box)
    else:
        print("[-] (init|leader|client|manager)")

if __name__ == "__main__":
    main()
