# register 
[![Build Status](https://travis-ci.org/docking-tools/register.svg?branch=master)](https://travis-ci.org/docking-tools/register)
[![GoDoc](https://godoc.org/github.com/docking-tools/register?status.svg)](https://godoc.org/github.com/docking-tools/register)

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
or use config file on folder $DOCKING_TOOLS/config.json

For consul, use http://<ip>:8500/v1/kv/<path>


## Templating
### env

Reads the given environment variable accessible to the current process.

{{env "CLUSTER_ID"}}
This function can be chained to manipulate the output:

{{env "CLUSTER_ID" | toLower}}

### convertGraphTopath

Read MetadataGraph and convert them to map key=URL_Path value
```
{{ range $key, $val := (convertGraphTopath .MetaDataGraph)}}{{$key}} {{$val}}\n{{end}}
```