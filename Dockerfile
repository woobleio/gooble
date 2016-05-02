FROM golang:latest

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
ENV GOROOT /usr/local/go
ENV PORT 8080

RUN mkdir -p $GOPATH/src/wobblapp && chmod -R 777 /go

WORKDIR $GOPATH/src/wobblapp

COPY . $GOPATH/src/wobblapp

RUN go get github.com/codegangsta/gin
RUN go-wrapper download
RUN go-wrapper install

EXPOSE 3000

CMD gin run
