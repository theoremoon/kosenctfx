import sys
import subprocess
import json
import requests
from pathlib import Path

def iterate_challenges(basedir):
    for taskjson in Path(basedir).glob("**/task.json"):
        with open(taskjson) as f:
            y = json.load(f)
        yield (taskjson.parent, y)


def main():
    basedir = Path(sys.argv[1])
    with open("challenges.json") as f:
        challenges = json.load(f)

    with open("server.json") as f:
        server = json.load(f)

    for c in iterate_challenges(basedir):
        cdir, cinfo = c
        if cinfo["name"] not in challenges:
            continue

        compose_path = cdir / "solution" / "docker-compose.yml"
        if not compose_path.is_file():
            continue

        if ("host" not in cinfo) or ("port" not in cinfo):
            continue

        try:
            print("[+] solve {}...".format(cinfo["name"]))
            r = subprocess.run(["docker-compose", "build"], timeout=300, cwd=compose_path.parent.as_posix())
            r = subprocess.run([
                "docker-compose",
                "run",
                "-e", "HOST={}".format(cinfo["host"]),
                "-e", "POST={}".format(cinfo["port"]),
                "solve"], timeout=300, cwd=compose_path.parent.as_posix(), stdout=subprocess.PIPE)
            if cinfo["flag"] in r.stdout.decode():
                solved = True
            else:
                solved = False
        except subprocess.TimeoutExpired:
            solved = False

        requests.post(
            server["url"] + "/admin/set-challenge-status",
            data=json.dumps({"name": cinfo["name"], "result": solved}),
            headers={
                "Content-Type": "application/json",
                "Authorization": "Bearer " + server["token"],
            }
        )


if __name__ == "__main__":
    main()
