FROM ubuntu
RUN apt-get update
RUN apt-get install -y wget
RUN wget https://go.dev/dl/go1.20.3.linux-amd64.tar.gz
RUN rm -rf /usr/local/go
RUN apt-get install tar -y 
RUN tar -C /usr/local -xzf go1.20.3.linux-amd64.tar.gz
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
ENV CGO_ENABLED=1

RUN apt-get  install -y libdlib-dev libblas-dev libatlas-base-dev liblapack-dev libjpeg-turbo8-dev gcc g++

COPY . .
RUN go mod download


EXPOSE 4000

ENTRYPOINT  go run ./cmd/web/*

