docker kill cl1 serv
docker rm cl1 serv

cd server
docker build -t server .
docker run -p 8080:8080 -d --name serv server

cd ../client
docker build -t client .
docker run  -i -t --link serv --name cl1 client