FROM scratch

#ADD . /go/src/github.com/docking-tools/register

#RUN go get github.com/tools/godep

ENV REGISTER_VERSION 0.0.7
add https://github.com/docking-tools/register/releases/download/0.0.7/register /register
#RUN cd /go/src/github.com/docking-tools/register && godep restore && go install
COPY example/config.json /config.json

#RUN chmod +x register && apt-get remove -y wget
#ADD
#ENV PATH $PATH:/register

ENTRYPOINT ["/register"]

CMD [""]