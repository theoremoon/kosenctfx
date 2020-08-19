rm -rf distfiles
mkdir distfiles
cp -r challenge/ distfiles/
rm distfiles/challenge/flag.txt
cp Dockerfile distfiles
cp docker-compose.yaml distfiles
