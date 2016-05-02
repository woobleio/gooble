FROM golang:latest

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
ENV GOROOT /usr/local/go

RUN mkdir -p $GOPATH/src/wobblapp && chmod -R 777 /go

ADD . $GOPATH/src/wobblapp

RUN cd $GOPATH/src/wobblapp && go install

ENTRYPOINT $GOPATH/bin/wobblapp

WORKDIR $GOPATH/src/wobblapp

EXPOSE 8000
