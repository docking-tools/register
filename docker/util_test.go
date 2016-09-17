package docker

import (
	"testing"
	"github.com/docker/engine-api/types/container"
)

func TestServiceMetadata(t *testing.T) {
	config := container.Config {
		Labels : make(map[string]string),
		Env : make([]string, 1),
	}

	config.Labels["service.test"]="ok"
	config.Labels["service.8080.test"]="ok-port"

	config.Labels["Service.8A.test"]="ko"
	config.Labels["service_test"]="ok"

	config.Env[0]="SERVICE_TEST=ok"

	metadata, metaFromPort :=  serviceMetaData(&config, "8080")

	t.Log("%v", metadata)

	if len(metadata)!=4 {
		t.Fatal("Number of result MetaData is not 4")
	}

	if len(metaFromPort) !=1 {
		t.Fatal("Number of result MetaFromPort is not 1")
	}	
}
