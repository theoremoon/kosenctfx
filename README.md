# KosenCTFx

CTF Scoreserver for InterKosenCTF 2020

## Files

```
.
├── ansible
│   └── docker.yml           -- ansible scripts for installing docker
├── docker-compose.yml
├── Dockerfile
├── envfile_example
├── Makefile
├── README.md
├── scoreserver              -- backend (golang)
├── solvability_check        -- toy script for solvability check
├── tool                     -- ctf management tool (python)
│   ├── ctf_config.yml
│   ├── config.ini_example
│   └── tool.py
└── ui                       -- frontend (vue)
```

## Usage (test / development)

up database / redis / file storate
```
$ docker-compose up
```

run scoreserver
```
$ make run
```

run frontend
```
$ cd ui; yarn serve
```

## Note

 This project is mainly for personal use. If you have some troubles and/or questions, tell the author by Issue.

## License

- TBD

## Author

- theoremoon
