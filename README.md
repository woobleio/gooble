# Wobble Backend

```
docker build -t gowobble .
docker run --name wobbleio -v "$PWD":/go/src/wobblapp -p 6060:8000 gowobble
```
