FROM circleci/golang:1.11

RUN wget https://packages.confluent.io/deb/5.1/archive.key 
RUN sudo apt-key add archive.key && rm archive.key
RUN sudo apt-get install apt-transport-https software-properties-common
RUN sudo add-apt-repository "deb http://deb.debian.org/debian stretch-backports main contrib non-free"
RUN sudo add-apt-repository "deb [arch=amd64] https://packages.confluent.io/deb/5.1 stable main"
RUN sudo apt-get install librdkafka-dev librdkafka1
