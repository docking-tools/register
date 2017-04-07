FROM alpine

ENV REGISTER_VERSION 0.0.7
add https://github.com/docking-tools/register/releases/download/0.0.7/register /register
COPY example/config.json  /root/.docking/config.json

RUN chmod +x /register

ENTRYPOINT ["/register"]
