import requests
import json
from string import ascii_lowercase
import random
import re
import time

URL = "http://localhost:5000"


def randomword():
    VOWELS = ['a', 'e', 'i', 'o', 'u']
    CONSONANTS = [ c for c in ascii_lowercase if c not in VOWELS] + ['']*2
    word = ''
    random_range = random.randint(3,4)
    for i in range(random_range):
        word += random.choice(VOWELS) + random.choice(CONSONANTS)
    return word

class Client():
    def __init__(self):
        self._cookies = None
        self._username = ""
        self._password = ""
        self._email = ""
        self._teamname = ""
        self._teamid = ""
        self._userid = ""
        self._challenges = []

    def _post(self, endpoint, data):
        r = requests.post(URL + endpoint, headers={"Content-Type": "application/json"}, data=json.dumps(data), cookies=self._cookies)

        print("==== REQUEST  =====")
        print(" URL: {}".format(endpoint))
        print(" DATA: {}".format(json.dumps(data, ensure_ascii=False)))
        print("")
        print("==== RESPONSE =====")
        print(" {}".format(r.text))
        print("")

        return r

    def _get(self, endpoint):
        r = requests.get(URL + endpoint, headers={"Content-Type": "application/json"}, cookies=self._cookies)

        print("==== REQUEST  =====")
        print(" URL: {}".format(endpoint))
        print("")
        print("==== RESPONSE =====")
        print(" {}".format(r.status_code))
        print(" {}".format(r.text))
        print("")
        return r


    def register_with_team(self):
        self._username = randomword()
        self._password = randomword()
        self._email = self._username + "@" + "example.com"
        self._teamname = randomword()

        self._post("/register-with-team", {
            "username": self._username,
            "teamname": self._teamname,
            "email": self._email,
            "password": self._password,
        })

    def login(self):
        r = self._post("/login", {
            "username": self._username if random.randint(0, 5) < 5 else randomword(),
            "password": self._password if random.randint(0, 5) < 5 else randomword(),
        })
        self._cookies = r.cookies

    def logout(self):
        r = self._post("/logout", {})
        self._cookies = r.cookies

    def info(self):
        r = self._get("/info")
        j = r.json()
        self._teamid = j.get("teamid", "")
        self._userid = j.get("userid", "")

    def info_update(self):
        r = self._get("/info-update")
        j = r.json()
        self._challenges = j.get("challenges", [])

    def renew_teamtoken(self):
        self._post("/renew-teamtoken", {})

    def getteam(self):
        self._get("/team/{}".format( self._teamid))

    def getuser(self):
        self._get("/user/{}" .format( self._userid))

    def submit(self):
        if len(self._challenges) == 0:
            self._post("/submit", {"flag": randomword()})
            return

        c = random.choice(self._challenges)
        flags = re.findall("KosenCTF{.+?}", c["description"])
        if len(flags) == 0:
            return
        flag = flags[0]
        
        self._post("/submit", {"flag": flag if random.randint(0, 5) < 4 else randomword()})
        self._get("/info-update?refresh=1")

client = Client()
methods = [f for f in dir(client) if f[0] != "_"]
client.register_with_team()
client.login()

while True:
    m = random.choice(methods)
    getattr(client, m)()
    time.sleep(random.randint(0, 60))

