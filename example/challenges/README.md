# InterKosenCTF2020-prectf

- category/title/[challenge | distfiles | solution] に問題を入れる
- ブランチ名は問題ごとに分けて category/title にする
- category/title/task.json 的な感じで問題ごとに設定できたらconflictしない&見やすいので嬉C

## category/title/task.json

```json
{
  "name": "challenge_name",
  "description": "challenge description. it is html",
  "flag": "KosenCTF{some_awesome_flag_wowow_takoyaki}",
  "author": "author name",
  "tags": ["crypto", "warmup"],
  "host": "pwn.kosenctf.com",
  "port": 8080,
  "is_survey": false
}
```

host / port は省略可


## Dockerfile
動かすプログラムに応じて次のコンテナを優先して使う。（上にくるほど優先度高。他にあったら追記して。）

- python:3.7-alpine
- php:7.4.1-fpm-alpine
- ubuntu:18.04

問題設定的にどうしても他のバージョンが必要な場合は他のコンテナを使っても良い。
（例:Flask問だが`LD_PRELOAD`などのハックが必要なのでubuntuコンテナを使う、など。）
