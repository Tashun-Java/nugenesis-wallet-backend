
docker-compose up --build

go run main.go

sudo docker build -f ./deploy/Dockerfile -t backend:2.0 .
sudo docker save -o backend-v2.tar backend:2.0

sudo docker images
sudo docker load -i backend-v2.tar 
sudo docker rm -f 3e0642e3ef3e
sudo docker rmi -f 3e0642e3ef3e

docker run -d -p 8080:8080 backend:2.0
