docker stop gg
docker rm gg
docker rmi team142/gg-mapgen:local
docker build -t team142/gg-mapgen:local .
docker run --name gg-mapgen --publish 8081:8081 team142/gg-mapgen:local -e "GG_MAP_PATH=/map/" -e "GG_MAP_LISTEN=:8081"