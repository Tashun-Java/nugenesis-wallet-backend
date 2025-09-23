
Ways to Checkout 
```bash

    git clone git@github.com:Tashun-Java/nugenesis-wallet-backend.git
    git submodule init
    git submodule update
```
Or 
```bash
  git clone --recurse-submodules git@github.com:Tashun-Java/nugenesis-wallet-backend.git

```


To Run the Backend

```bash
  go run ./cmd/app/main.go
```


docker-compose up --build


Docker Build and Run with helpful docker commands 
```bash

sudo docker build -f ./deploy/Dockerfile -t backend:2.0 .
sudo docker save -o backend-v2.tar backend:2.0

sudo docker images
sudo docker load -i backend-v2.tar
sudo docker rm -f 3e0642e3ef3e
sudo docker rmi -f 3e0642e3ef3e

docker run -d -p 8080:8080 backend:2.0
```

