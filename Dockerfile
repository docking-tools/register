FROM debian

ENV REGISTRY_VERSION 0.0.1
ENV BIN_URL https://github.com/docking-tools/register/releases/download/0.0.1/register

COPY example/config.json /register/config.json

RUN cd /register && apt-get update && apt-get install -y wget && wget $BIN_URL && chmod +x register && apt-get remove -y wget
ENV PATH $PATH:/register

COPY entrypoint.sh /
ENTRYPOINT ["/entrypoint.sh"]

CMD ["register"]