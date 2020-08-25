# KosenCTFx

## 周辺環境系

主に`tool/tool.py`を用いるところと`infra/infra.py`を用いるところがある。

`infra/infra.py`は`init` みたいなサブコマンドが用意されているけどまだ安定してない

### Bucket の初期化

`python ./tool.py bucket_init` とする

### Managerの初期化

`python ./tool.py manager_init` とする

## CTFの管理系

主に`tool/tool.py`を用いる。先に`config.ini`をいい感じに設定しておく必要がある。
くわしくは

### ScoreServerの設定（CTFの開催時刻など）を変更する

`ctf_config.yml`を設定した上で、 `python ./tool.py apply_config ../ctf_config.yml`

### 問題の登録

問題を取り扱う場合は`python ./tool.py --basedir <challenges dir path>` から始める

`python ./tool.py --basedir <challenges dir path> register`

### 問題のオープン

目当ての問題のIDを `python ./tool.py --basedir <challenges dir path> list | jq '...'`などとして手に入れておき
`python ./tool.py --basedir <challenges dir path> open <id>`


### その他

ヘルプを見てくれ

