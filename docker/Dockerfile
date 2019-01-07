FROM buildpack-deps:xenial-scm

# gcc for cgo
RUN apt-get update && apt-get install -y --no-install-recommends \
		g++ \
		gcc \
		libc6-dev \
		make \
		pkg-config \
        apt-utils \
        apt-transport-https	

ENV GOLANG_VERSION 1.11.4

RUN set -eux; \
	\
	url="https://golang.org/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz"; \
	wget -O go.tgz "$url"; \
	tar -C /usr/local -xzf go.tgz; \
	rm go.tgz; \
	export PATH="/usr/local/go/bin:$PATH"; \
	go version

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH

ENV DEBIAN_FRONTEND noninteractive

# Install rdkafka and deps
RUN wget https://packages.confluent.io/deb/5.1/archive.key && \
    apt-key add archive.key && rm archive.key && \
    echo "deb [arch=amd64] https://packages.confluent.io/deb/5.1 stable main" | tee -a /etc/apt/sources.list && \
    apt-get update && apt-get install -y --no-install-recommends librdkafka-dev librdkafka1 librdkafka++1 && \
    rm -rf /var/lib/apt/lists/*
