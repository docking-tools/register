# register
This tools use  [Golang template](https://golang.org/pkg/text/template/)
## Usage
```sh
export <docker_event_uppercase>_TMPL_<tempalte_name>=<HTTP_CMD>:/my-query data
export DOCKER_HOST=tcp://<ip:port>

# Run listener
register http://<ip:port>
```