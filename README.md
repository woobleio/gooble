# Wobble Backend
AY LMAO!

## Installation

```
docker build -t gowobble .
docker run --name wobbleio -v "$PWD":/go/src/wobblapp -p 3000:3000 -d gowobble
```
## Configuration

```
docker exec -it wobbleio /bin/bash
cd $GOPATH/configs
echo "db_host:<YOUR_MONGODB_HOST>" >> config.yml
echo "db_port:<YOUR_MONGODB_PORT>" >> config.yml
echo "db_username:<YOUR_MONGODB_USER>" >> config.yml
echo "db_password:<YOUR_USER_PASSWD>" >> config.yml
```
