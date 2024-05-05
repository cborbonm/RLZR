# Use base Linux distribution image
FROM ubuntu:latest

# Switch to root mode
USER root

# Update package repository and install git
RUN apt-get update && apt-get install -y git vim

# Install GoLang
RUN apt-get install -y golang

# Install ZMap Dependencies
RUN apt-get install -y build-essential cmake libgmp3-dev gengetopt libpcap-dev flex byacc libjson-c-dev pkg-config libunistring-dev libjudy-dev iptables

# Install ZMap
RUN apt-get install -y zmap

# Set environment variables for Go
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

# Clone rlzr repo
RUN git clone https://github.com/cborbonm/RLZR.git rlzr
RUN cd rlzr && git checkout setup && cd ..

# Start Go Repo
RUN go mod init rlzr

# Download rlzr Go dependencies
RUN go get github.com/stanford-esrg/lzr/bin
RUN go get github.com/stanford-esrg/lzr/handshakes

# Set the working directory inside the container
WORKDIR /rlzr

# Set an entry point or default command if needed
ENTRYPOINT ["/bin/bash"]



#FROM golang:1.14

#RUN apt-get update && apt-get install -y libpcap-dev

#RUN go get -v gopkg.in/mgo.v2/bson
#RUN go get -v github.com/stanford-esrg/lzr

#COPY . /go/src/github.com/stanford-esrg/lzr/

#RUN (cd /go/src/github.com/stanford-esrg/lzr && make lzr)

#WORKDIR /go/src/github.com/stanford-esrg/lzr

#CMD ["lzr"]
