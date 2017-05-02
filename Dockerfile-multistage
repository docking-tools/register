FROM golang:1.7.3 as builder

WORKDIR /go/src/github.com/docking-tools/register/

RUN wget "https://github.com/Masterminds/glide/releases/download/0.10.2/glide-0.10.2-linux-amd64.tar.gz"  && \
  mkdir -p $HOME/bin && tar -vxz -C $HOME/bin --strip=1 -f glide-0.10.2-linux-amd64.tar.gz && \
  export PATH="$HOME/bin:$PATH" && \
  glide install --strip-vendor --strip-vcs

COPY . .  

RUN go test -v -cover $(glide novendor) && \
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build && \
  chmod register
  
  
FROM alpine:latest

WORKDIR /root/

ENV REGISTER_VERSION 0.0.7
COPY --from=builder /go/src/github.com/docking-tools/register/register /register
COPY --from=builder example/config.json  /root/.docking/config.json

CMD ["/register"]
