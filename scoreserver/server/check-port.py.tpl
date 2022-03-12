import socket
import threading
import datetime

JST = datetime.timezone(datetime.timedelta(hours=+9), 'JST')

START = datetime.datetime(2021, 12, 11, 14)
challs = [
    {{- range $idx, $c := .Challenges }}
    {{ if $c.Host }}
    {{- "" -}}{"name": "{{- $c.Name -}}", "port": {{ $c.Port }},  "host": "{{- $c.Host -}}", "release": START},
    {{ end -}}
    {{ end }}
]

def check(chall):
    name, port, host, release = chall['name'], chall['port'], chall['host'], chall['release']
    release = release.astimezone(JST)
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.settimeout(1)
    try:
        r = s.connect_ex((host, port))
        if r == 0:
            print(f"\x1b[32m\x1b[1m{name}: OPEN ({host}:{port})\x1b[0m", end='')
            if release <= datetime.datetime.now(JST):
                print(f"\t--> \x1b[36m\x1b[1mOK (Schedule: {release})\x1b[0m")
            else:
                print(f"\t--> \x1b[31m\x1b[1mClose It (Schedule: {release})\x1b[0m")
        else:
            print(f"\x1b[33m\x1b[1m {name}: CLOSED ({host}:{port})\x1b[0m", end='')
            if release > datetime.datetime.now(JST):
                print(f"\t--> \x1b[36m\x1b[1mOK (Schedule: {release})\x1b[0m")
            else:
                print(f"\t--> \x1b[31m\x1b[1mOpen It (Schedule: {release})\x1b[0m")
        s.close()
    except socket.gaierror as e:
        print(f"\x1b[41m\x1b[37m\x1b[1m {name}: Critical ({e})\x1b[0m")

if __name__ == '__main__':
    thList = []
    for chall in challs:
        th = threading.Thread(target=check, args=(chall,))
        th.start()
        thList.append(th)
    for th in thList:
        th.join()
