git pull origin master &&
docker container stop wealthy-backend &&
docker container rm wealthy-backend &&
docker image rm wealthy-backend &&
docker build . -t wealthy-backend &&
docker run -d --name wealthy-backend -p 8888:8080 wealthy-backend
