# Wobble Backend

```
docker build -t gowobble .
docker run --name wobbleio -v "$PWD":/go/src/wobblapp -p 3000:3000 -d gowobble
```
