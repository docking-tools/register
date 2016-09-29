package docker

import (
	"testing"
	"github.com/docker/engine-api/types/container"
	"github.com/docking-tools/register/api"
	"strings"
)

func TestServiceMetadataSSS(t *testing.T) {
	config := container.Config {
		Labels : make(map[string]string),
		Env : make([]string, 1),
	}

	config.Labels["SERVICE.NAME"]="test1"
	config.Labels["service_8080_name"]="ok-port"
	config.Labels["SERVICE_ignore"]="true"

	config.Labels["service.8A_test"]="ko"

	config.Labels["service_test"]="ok"
	config.Labels["test_service_test"]="ko"

	config.Env[0]="SERVICE.TEST=ok"

	metadata, metaFromPort :=  serviceMetaData(&config, "8080")

	t.Log("%v", metadata)
	t.Log("%v", metaFromPort)

	ignore := mapDefault(metadata,"ignore","")
	t.Log("%#v", ignore)


	if len(metadata)!=4 {
		t.Fatal("Number of result MetaData is not 4")
	}
	if !metaFromPort["name"] {
		t.Fatal("mettaFromPort for key name can be true")
	}
	if !strings.EqualFold(metadata["ignore"],"true") {
		t.Fatal("mettadata for key 'ignore' can be true")
	}


	if len(metaFromPort) !=1 {
		t.Fatal("Number of result MetaFromPort is not 1")
	}
}

func TestGraphMetaData(t *testing.T) {
	config := container.Config {
		Labels : make(map[string]string),
		Env : make([]string, 1),
	}

	config.Labels["cron.test.titi"]="ok"
	config.Labels["cron.test"]="KO"

	config.Labels["cron.test.tutu"]="ok"
	config.Labels["cron.8080.test"]="ok-port"

	config.Labels["crone.8A.test"]="ko"
	config.Labels["cron_test.toto.tata"]="ok"
	config.Env[0]="test_cron=ok"


	result  :=  graphMetaData(&config)
	t.Logf("%v", result)

	if len(result)!=3 {
		t.Fatal("Number of result MetaData is not 3 %v", result)
	}

	if len(result["cron"].(api.Recmap)) !=2 {
		t.Fatal("Number of result MEtaData is not 1 %v", result["cron"].(api.Recmap))
	}
	if result["test"].(api.Recmap)["cron"]!="ok" {
		t.Fatal("cron.test not equals to ok", result["cron"].(api.Recmap)["test"])
	}
	if result["cron"].(api.Recmap)["test"].(api.Recmap)["tutu"]!="ok" {
		t.Fatal("cron.test.tutu not equals to ok", result["cron"].(api.Recmap)["test"])
	}
}