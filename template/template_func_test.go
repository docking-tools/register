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
	entry["c"]=1

	resul := convertGraphTopath(entry)
	if len(resul)==3 {
		t.Error("number of result not equals to 2 %v", resul)
	}
	if resul["/a/b"]!="2" {
		t.Error("key /a/b not equals to 2 %v", resul)
	}
	if resul["/a/c"]!="3" {
		t.Error("key /a/c not equals to 3 %v", resul)
	}
}