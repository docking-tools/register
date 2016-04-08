FROM busybox

ENV REGISTRY_VERSION 0.0.1
ENV BIN_URL https://github.com/docking-tools/register/releases/download/0.0.1/register


COPY $BIN_URL /register/register
COPY example/config.json /register/config.json

CMD register -c /register/config.json