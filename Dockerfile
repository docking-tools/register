FROM golang:latest

ADD . /go/src/github.com/docking-tools/register

RUN go get github.com/tools/godep

#ENV REGISTER_VERSION 0.0.4
#ENV BIN_URL https://github.com/docking-tools/register/releases/download/0.0.4/register
RUN cd /go/src/github.com/docking-tools/register && godep restore && go install
COPY example/config.json /config.json

#RUN cd /register && apt-get update && apt-get install -y wget && wget $BIN_URL && chmod +x register && apt-get remove -y wget
#ADD
#ENV PATH $PATH:/register

COPY entrypoint.sh /
ENTRYPOINT ["/entrypoint.sh"]

CMD ["register"]