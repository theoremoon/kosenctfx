import requests
import random
import urllib.parse
from string import ascii_lowercase
from pathlib import Path
import os
import json

def challengename():
    def get_word():
        prefix_len = len("https://en.wikipedia.org/wiki/")
        if random.randint(1, 2) == 1:
            r = requests.get("https://en.wikipedia.org/wiki/Special:Random")
        else:
            r = requests.get("https://jp.wikipedia.org/wiki/Special:Random")
        return urllib.parse.unquote(r.url[prefix_len:])

    return "{} {}".format(get_word(), get_word())

def challengeid():
    VOWELS = ['a', 'e', 'i', 'o', 'u']
    CONSONANTS = [ c for c in ascii_lowercase if c not in VOWELS] + ['']*2
    word = ''
    random_range = random.randint(3,4)
    for i in range(random_range):
        word += random.choice(VOWELS) + random.choice(CONSONANTS)
    return word

def challengetag():
    tags = ["crypto", "web", "pwn", "reversing", "misc", "warmup"]
    return random.sample(tags, k=random.randint(0, 3))

def author():
    return random.choice(["theoremoon", "yoshiking", "ptr-yudai", "insecure"])

def get_random_unicode(length):
    try:
        get_char = unichr
    except NameError:
        get_char = chr

    # Update this to include code point ranges to be sampled
    include_ranges = [
        ( 0x0021, 0x0021 ),
        ( 0x0023, 0x0026 ),
        ( 0x0028, 0x007E ),
        ( 0x00A1, 0x00AC ),
        ( 0x00AE, 0x00FF ),
        ( 0x0100, 0x017F ),
        ( 0x0180, 0x024F ),
        ( 0x2C60, 0x2C7F ),
        ( 0x16A0, 0x16F0 ),
        ( 0x0370, 0x0377 ),
        ( 0x037A, 0x037E ),
        ( 0x0384, 0x038A ),
        ( 0x038C, 0x038C ),
    ]

    alphabet = [
        get_char(code_point) for current_range in include_ranges
            for code_point in range(current_range[0], current_range[1] + 1)
    ]
    while True:
        r = ''.join(random.choice(alphabet) for i in range(length))
        if "$" not in r:
            return r
    

dir = Path(__file__).parent / "challenges" / "dummy"
id = challengeid()
os.makedirs(dir / id)
with open(dir / id / "task.json", "w") as f:
    json.dump({
        "name": challengename(),
        "description": "the flag is <pre>${flag}</pre>" +  get_random_unicode(random.randrange(1000)),
        "flag": "KosenCTF{{{}}}".format(id),
        "author": author(),
        "tags": challengetag(),
        "is_survey": False,
    }, f)
