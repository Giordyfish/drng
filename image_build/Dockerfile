FROM golang:1.18

ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt upgrade -y
RUN apt-get install git
RUN git clone https://github.com/Giordyfish/drng
WORKDIR /go/drng
RUN go mod vendor
RUN go install