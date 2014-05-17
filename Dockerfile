FROM debian:jessie

RUN apt-get update && apt-get install --no-install-recommends -y \
    ca-certificates \
    curl \
    mercurial \
    git-core

RUN curl -s https://storage.googleapis.com/golang/go1.2.2.linux-amd64.tar.gz | tar -v -C /usr/local -xz

ENV GOPATH /go
ENV GOROOT /usr/local/go
ENV PATH $PATH:/usr/local/go/bin:/go/bin

RUN git clone https://github.com/kairichard/zreqlo.git app
RUN cd app; go get; go build

EXPOSE 8080

CMD app/zreqlo -redis="$REDIS_PORT_6379_TCP_ADDR:$REDIS_PORT_6379_TCP_PORT" -bind="0.0.0.0:8080"
