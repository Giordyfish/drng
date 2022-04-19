FROM golang:alpine

ARG DEBIAN_FRONTEND=noninteractive

RUN apk update
RUN apk add git
RUN git clone https://github.com/Giordyfish/drng
WORKDIR /go/drng
RUN go mod vendor
RUN go install
ENTRYPOINT ["/go/drng/demo/data/client-script.sh"]