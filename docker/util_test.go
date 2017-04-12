package docker

import (
	"github.com/docker/docker/api/types"
	container "github.com/docker/docker/api/types/container"
	swarm "github.com/docker/docker/api/types/swarm"
	"github.com/docking-tools/register/api"
	"strings"
	"testing"
)

func TestServiceMetadataSSS(t *testing.T) {
	config := container.Config{
		Labels: make(map[string]string),
		Env:    make([]string, 1),
	}

	config.Labels["SERVICE.NAME"] = "test1"
	config.Labels["service_8080_name"] = "ok-port"
	config.Labels["SERVICE_ignore"] = "true"

	config.Labels["service.8A_test"] = "ko"

	config.Labels["service_test"] = "ok"
	config.Labels["test_service_test"] = "ko"

	config.Env[0] = "SERVICE.TEST=ok"

	metadata, metaFromPort := serviceMetaData("8080", envArrayToMap(config.Env), config.Labels)

	t.Log("%v", metadata)
	t.Log("%v", metaFromPort)

	ignore := mapDefault(metadata, "ignore", "")
	t.Log("%#v", ignore)

	if len(metadata) != 4 {
		t.Fatal("Number of result MetaData is not 4")
	}
	if !metaFromPort["name"] {
		t.Fatal("mettaFromPort for key name can be true")
	}
	if !strings.EqualFold(metadata["ignore"], "true") {
		t.Fatal("mettadata for key 'ignore' can be true")
	}

	if len(metaFromPort) != 1 {
		t.Fatal("Number of result MetaFromPort is not 1")
	}
}

func TestGraphMetaData(t *testing.T) {
	config := container.Config{
		Labels: make(map[string]string),
		Env:    make([]string, 2),
	}

	config.Labels["cron.test.titi"] = "ok"
	//	config.Labels["cron.test"] = "KO"

	config.Labels["cron.test.tutu"] = "ok"
	config.Labels["cron.8080.test"] = "ok-port"

	config.Labels["crone.8A.test"] = "ko"
	config.Labels["cron_test.toto.tata"] = "ok"
	config.Env[0] = "test_cron=ok"
	config.Env[1] = "sans_valeur"

	result := graphMetaData(envArrayToMap(config.Env), config.Labels)
	t.Logf("%v", result)

	if len(result) != 3 {
		t.Fatal("Number of result MetaData is not 3 %v", result)
	}

	if len(result["cron"].(api.Recmap)) != 2 {
		t.Fatal("Number of result MEtaData is not 1 %v", result["cron"].(api.Recmap))
	}
	if result["test"].(api.Recmap)["cron"] != "ok" {
		t.Fatal("cron.test not equals to ok", result["cron"].(api.Recmap)["test"])
	}
	if result["cron"].(api.Recmap)["test"].(api.Recmap)["tutu"] != "ok" {
		t.Fatal("cron.test.tutu not equals to ok", result["cron"].(api.Recmap)["test"])
	}
}

func TestServiceSwarmPort(t *testing.T) {
	container := types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			ID: "idddd",
		},
		Config: &container.Config{
			Hostname: "rrr",
		},
	}

	swarmService := swarm.Service{
		Spec: swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Name: "serviceTest",
			},
		},
	}

	port := swarm.PortConfig{
		TargetPort:    80,
		PublishedPort: 8080,
		PublishMode:   "tcp",
	}

	// call method
	result := serviceSwarmPort(&container, &swarmService, port)

	if result.HostPort != "8080" {
		t.Fatal("Host port not equals to 8080")
	}
	if result.ExposedPort != "80" {
		t.Fatal("Exposed port not equals to 80")
	}
	if result.PortType != "tcp" {
		t.Fatal("Port type not equals to tcp")
	}
}
