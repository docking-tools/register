FROM alpine

ENV REGISTER_VERSION test
ADD https://github.com/docking-tools/register/releases/download/test/register /register
COPY example/config.json  /root/.docking/config.json

RUN chmod +x /register

ENTRYPOINT ["/register"]