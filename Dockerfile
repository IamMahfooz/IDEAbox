
FROM ubuntu:latest
WORKDIR /ideabox

RUN apt-get update
RUN apt-get install -y libdlib-dev libblas-dev libatlas-base-dev liblapack-dev libjpeg-turbo8-dev gcc g++ wget tar

RUN wget https://go.dev/dl/go1.20.3.linux-amd64.tar.gz
RUN rm -rf /usr/local/go
RUN tar -C /usr/local -xzf go1.20.3.linux-amd64.tar.gz
RUN rm go1.20.3.linux-amd64.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

COPY . .
RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux go build -o bin/app cmd/web/*.go

ENTRYPOINT ["bin/app"]
