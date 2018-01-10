docker stop gg-mapgen
docker rm gg-mapgen
docker rmi team142/gg-mapgen:local
docker build -t team142/gg-mapgen:local .
docker run -d --name gg-mapgen --publish 8081:8081 team142/gg-mapgen:local -e GG_MAP_PATH=/ -e GG_MAP_LISTEN=0.0.0.0:8081