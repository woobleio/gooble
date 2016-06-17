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
      && touch $CONFPATH/dev.yml \
      && touch $CONFPATH/test.yml

#
# SET DATABASE CONF
#
RUN echo "db_host: wobbleservice_wobbledb_1" >> $CONFPATH/dev.yml \
      && echo "db_name: wobbleapp" >> $CONFPATH/dev.yml \
      && echo "db_port:" >> $CONFPATH/dev.yml \
      && echo "db_username:" >> $CONFPATH/dev.yml \
      && echo "db_password:" >> $CONFPATH/dev.yml

RUN echo "db_host: wobbleservice_wobbledb_1" >> $CONFPATH/test.yml \
      && echo "db_name: wobbleapp_test" >> $CONFPATH/test.yml \
      && echo "db_port:" >> $CONFPATH/test.yml \
      && echo "db_username:" >> $CONFPATH/test.yml \
      && echo "db_password:" >> $CONFPATH/test.yml

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
