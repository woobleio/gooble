FROM golang:latest

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
ENV GOROOT /usr/local/go

RUN mkdir -p $GOPATH/src/github.com/wobbleio && chmod -R 777 /go

ADD . $GOPATH/src/github.com/wobbleio

RUN cd $GOPATH/src/github.com/wobbleio && go install

ENTRYPOINT $GOPATH/bin/wobbleio

WORKDIR $GOPATH/src/github.com/wobbleio

EXPOSE 8080
