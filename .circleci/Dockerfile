FROM circleci/golang:1.11

RUN sudo apt-get install apt-utils apt-transport-https 
RUN wget https://packages.confluent.io/deb/5.1/archive.key && \
    sudo apt-key add archive.key && rm archive.key
RUN wget http://security-cdn.debian.org/debian-security/pool/updates/main/o/openssl/libssl1.0.0_1.0.1t-1+deb8u10_amd64.deb && \
    sudo dpkg -i libssl1.0.0_1.0.1t-1+deb8u10_amd64.deb && \
    rm libssl1.0.0_1.0.1t-1+deb8u10_amd64.deb
RUN echo "deb [arch=amd64] https://packages.confluent.io/deb/5.1 stable main" | sudo tee -a /etc/apt/sources.list && \
    sudo apt-get update && sudo apt-get install librdkafka-dev librdkafka1