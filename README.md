# register 
[![Build Status](https://travis-ci.org/docking-tools/register.svg?branch=master)](https://travis-ci.org/docking-tools/register)
This tools use  [Golang template](https://golang.org/pkg/text/template/)
## Usage
```sh
export <docker_event_uppercase>_TMPL_<tempalte_name>=<HTTP_CMD>:/my-query data
export DOCKER_HOST=tcp://<ip:port>

# Run listener
register http://<ip:port>
```
## Docker usage
```sh
docker run -it --rm -e HOST_IP="<public_ip>"" -e REGISTER_URL="http://xx.xx.xx.xx:xxxx" -e DOCKER_URL="" dockingtools/register:latest 
```
For consul, use http://<ip>:8500/v1/kv/<path>

