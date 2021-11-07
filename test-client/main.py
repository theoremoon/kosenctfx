import requests
import random
import string
import json
import glob
import yaml
import sys
import os
import time
import concurrent.futures
from typing import List

def random_str(n):
    buf = []
    for i in range(n):
        c = chr(0)
        while not c.isprintable():
            # c = chr(random.randrange(0x20, 0x110000))
            c = chr(random.randrange(0x20, 0x7f))
        buf.append(c)
    return "".join(buf)

def random_ascii(n):
    return "".join(random.choices(string.printable[:-10], k=n))

def randstr(n, k):
    if random.random() < k:
        return random_ascii(n)
    else:
        return random_str(n)

def random_country():
    country_list = "AF AX AL DZ AS AD AO AI AQ AG AR AM AW AU AT AZ BH BS BD BB BY BE BZ BJ BM BT BO BQ BA BW BV BR IO BN BG BF BI KH CM CA CV KY CF TD CL CN CX CC CO KM CG CD CK CR CI HR CU CW CY CZ DK DJ DM DO EC EG SV GQ ER EE ET FK FO FJ FI FR GF PF TF GA GM GE DE GH GI GR GL GD GP GU GT GG GN GW GY HT HM VA HN HK HU IS IN ID IR IQ IE IM IL IT JM JP JE JO KZ KE KI KP KR KW KG LA LV LB LS LR LY LI LT LU MO MK MG MW MY MV ML MT MH MQ MR MU YT MX FM MD MC MN ME MS MA MZ MM NA NR NP NL NC NZ NI NE NG NU NF MP NO OM PK PW PS PA PG PY PE PH PN PL PT PR QA RE RO RU RW BL SH KN LC MF PM VC WS SM ST SA SN RS SC SL SG SX SK SI SB SO ZA GS SS ES LK SD SR SJ SZ SE CH SY TW TJ TZ TH TL TG TK TO TT TN TR TM TC TV UG UA AE GB US UM UY UZ VU VE VN VG VI WF EH YE ZM ZW".split() + ([" "] * 10)
    return random.choice(country_list)

class Client:
    def __init__(self, base_url: str, flags: List[str]):
        self.session = requests.Session()
        self.base_url = base_url
        self.username = None
        self.password = None
        self.country = ''
        self.flags = flags + [flag.lower() for flag in flags[:10]] + [random_str(i) for i in range(0, 10)]

    def run(self, i):
        actions = [name for name in dir(self) if callable(getattr(self, name)) and not name.startswith("_") and name != "run"]

        time.sleep(random.randint(0, 60))
        self.register()
        time.sleep(random.randint(0, 60))
        self.login()
        while True:
            time.sleep(random.randint(0, 60))
            print(i)
            action = random.choice(actions)
            getattr(self, action)()

    def register(self):
        if self.username:
            return

        self.username = random_str(random.randint(1, 30))
        self.password = random_str(random.randint(1, 30))
        self.email = random_ascii(10) + "@example.com"
        self.country = random_country()
        return self._post("/api/register", {
            "teamname": self.username,
            "password": self.password,
            "email": self.email,
            "country": self.country,
        })

    def login(self):
        if self.username is None:
            return

        self._post("/api/login", {
            "teamname": self.username,
            "password": self.password,
        })

    def logout(self):
        self._post("/api/logout", {})

    def update_profile(self):
        username = random_str(random.randint(1, 30))
        password = random_str(random.randint(1, 30))
        country = random_country()
        r = self._post("/api/update-profile", {
            "teamname": self.username,
            "password": self.password,
            "country": self.country,
        })
        if 200 <= r.status_code < 300:
            self.username = username
            self.password = password
            self.country = country

    def getInfo(self):
        self._get("/api/info")
        self._get("/api/info-update")

    def submit(self):
        self._post("/api/submit", {
            "flag": random.choice(self.flags),
        })

    def _get(self, endpoint):
        return self.session.get(self.base_url + endpoint)

    def _post(self, endpoint, data):
        return self.session.post(self.base_url + endpoint, data=json.dumps(data), headers={"Content-Type": "application/json"})

def main():
    baseurl = sys.argv[1]
    n = int(sys.argv[2])
    taskdir = sys.argv[3]
    flags = []
    for y in glob.glob(os.path.join(taskdir, "**/task.yml"), recursive=True):
        with open(y) as f:
            flags.append(yaml.load(f.read())["flag"])

    with concurrent.futures.ThreadPoolExecutor(max_workers=n) as ex:
        futures = []
        for i in range(n):
            futures.append(ex.submit(lambda: Client(baseurl, flags).run(i)))

        concurrent.futures.wait(futures)


if __name__ == "__main__":
    main()
