FROM busybox

ENV REGISTRY_VERSION 0.0.1
ENV BIN_URL https://github.com/docking-tools/register/releases/download/0.0.1/register

COPY example/config.json /register/config.json

RUN cd /register && wget $BIN_URL

CMD register -c /register/config.json