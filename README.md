# Wobble Backend
AYY LMAO!

## Installation

```
docker build -t gowobble .
docker run --name wobbleio -v "$PWD":/go/src/wobblapp -p 3000:3000 -d gowobble
```
## Configuration

```
docker exec -it wobbleio /bin/bash
cd $GOPATH/configs
echo "db_host:<YOUR_MONGODB_HOST>" >> dev.yml
echo "db_name:<YOUR_MONGODB_NAME>" >> dev.yml
echo "db_port:<YOUR_MONGODB_PORT>" >> dev.yml
echo "db_username:<YOUR_MONGODB_USER>" >> dev.yml
echo "db_password:<YOUR_USER_PASSWD>" >> dev.yml
```
