FROM golang:latest

ENV GOENV dev
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
ENV GOROOT /usr/local/go
ENV PORT 8080
ENV CONFPATH $GOPATH/configs

RUN mkdir -p $GOPATH/src/wobblapp \
      && chmod -R 777 /go \
      && mkdir $CONFPATH \
      && touch $CONFPATH/dev.yml

WORKDIR $GOPATH/src/wobblapp

COPY . $GOPATH/src/wobblapp

RUN go get github.com/codegangsta/gin
RUN go get github.com/gin-gonic/gin
RUN go get github.com/smartystreets/goconvey
RUN go get github.com/spf13/viper
RUN go get gopkg.in/mgo.v2
RUN go get gopkg.in/mgo.v2/bson
RUN go-wrapper download
RUN go-wrapper install

EXPOSE 3000

CMD gin run
