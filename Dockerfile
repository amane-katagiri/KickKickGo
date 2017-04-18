FROM golang:alpine
MAINTAINER Amane Katagiri
CMD [""]
ENTRYPOINT ["kick-kick-go"]
WORKDIR /go/src/github.com/amane-katagiri/kick-kick-go
COPY glide.yaml Makefile /go/src/github.com/amane-katagiri/kick-kick-go/
RUN apk add --update make git && make deps
COPY example /app
COPY . /go/src/github.com/amane-katagiri/kick-kick-go/
RUN make && make install
WORKDIR /app
