package template

import (
    "os"
    "testing"
	"github.com/docking-tools/register/api"
)

func TestEnv(t *testing.T) {
	if err := os.Setenv("foo", "bar"); err != nil {
		t.Fatal(err)
	}

	result, err := env("foo")
	if err != nil {
		t.Fatal(err)
	}

	if result != "bar" {
		t.Errorf("expected %#v to be %#v", result, "bar")
	}
}

func TestConvertGraphTopath(t *testing.T) {
	entry := make(api.Recmap, 0)
	entry["a"] = make(api.Recmap, 0)
	entry["a"].(api.Recmap)["b"]="2"
	entry["a"].(api.Recmap)["c"]="3"
	entry["a"].(api.Recmap)["d"]=make(api.Recmap, 0)
	entry["a"].(api.Recmap)["d"].(api.Recmap)["e"]="4"
	entry["c"]=1

	entry["cron"]=make(api.Recmap, 0)
	entry["cron"].(api.Recmap)["test"]=make(api.Recmap, 0)
	entry["cron"].(api.Recmap)["test"].(api.Recmap)["cmd"]="ps"

	t.Logf("entry %v", entry)

	resul := convertGraphTopath(entry)
	if len(resul)!=4 {
		t.Errorf("number of result not equals to 4 / %v %v", len(resul), resul)
	}
	if val,ok := resul["/a/b"];ok && val!="2" {
		t.Errorf("key /a/b not equals to 2 %v", resul)
	}
	if resul["/a/c"]!="3" {
		t.Errorf("key /a/c not equals to 3 %v", resul)
	}
}

func TestListPathfromGraph(t *testing.T) {
	entry := make(api.Recmap, 0)
	entry["a"] = make(api.Recmap, 0)
	entry["a"].(api.Recmap)["b"]="2"
	entry["a"].(api.Recmap)["c"]="3"
	entry["a"].(api.Recmap)["d"]=make(api.Recmap, 0)
	entry["a"].(api.Recmap)["d"].(api.Recmap)["e"]="4"
	entry["c"]=1

	resul := listPathfromGraph(entry)
	t.Log(resul)
	if len(resul)==3 {
		t.Error("number of result not equals to 2 %v", resul)
	}
	if !stringInSlice("/a", resul) {
		t.Error("key /a not found %v", resul)
	}
	if !stringInSlice("/a/d", resul) {
		t.Error("key /a/d not found %v", resul)
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}