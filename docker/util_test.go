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

	config.Labels["service.name"]="test1"
	config.Labels["service.8080.name"]="ok-port"
	config.Labels["service.ignore"]="true"

	config.Labels["Service.8A.test"]="ko"

	config.Labels["service_test"]="ok"
	config.Labels["test_service_test"]="ko"

	config.Env[0]="SERVICE_TEST=ok"

	metadata, metaFromPort :=  serviceMetaData(&config, "8080")

	t.Log("%v", metadata)
	t.Log("%v", metaFromPort)

	ignore := mapDefault(metadata,"ignore","")
	t.Log("%#v", ignore)
	if len(metadata)!=4 {
		t.Fatal("Number of result MetaData is not 4")
	}
	//if metaFromPort["8080"]

	if len(metaFromPort) !=1 {
		t.Fatal("Number of result MetaFromPort is not 1")
	}	
}

func TestGraphMetaData(t *testing.T) {
	config := container.Config {
		Labels : make(map[string]string),
		Env : make([]string, 1),
	}

	config.Labels["cron.test"]="ok"
	config.Labels["cron.8080.test"]="ok-port"

	config.Labels["crone.8A.test"]="ko"
	config.Labels["cron_test"]="ok"
	config.Labels["test_cron_test"]="ko"

	config.Env[0]="CRON_TEST=ok"

	result  :=  graphMetaData(&config, "cron")

	t.Log("%v", result)

	if len(result)!=1 {
		t.Fatal("Number of result MetaData is not 1")
	}

	if len(result["cron"].(recmap)) !=2 {
		t.Fatal("Number of result MetaFromPort is not 1")
	}	
}