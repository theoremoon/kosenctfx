# KosenCTFx

CTF Scoreserver for InterKosenCTF 2020

## Files

```
.
├── docker-compose.yml
├── docker-compose.prod.yml
├── Dockerfile
├── Makefile
├── README.md
├── scoreserver              -- backend (golang)
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
