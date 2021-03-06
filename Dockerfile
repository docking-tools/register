FROM alpine

ENV REGISTER_VERSION 0.0.9
ADD https://github.com/docking-tools/register/releases/download/0.0.9/register /register
COPY example/config.json  /root/.docking/config.json

RUN chmod +x /register

CMD ["/register"]