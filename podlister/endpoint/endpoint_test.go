package endpoint

import (
	"bytes"
	"io/ioutil"
	"testing"
)

type ReadFiler struct {
	path string
}

func (r ReadFiler) ReadFile(path string) ([]byte, error) {
	buf := bytes.NewBufferString(r.path)
	return ioutil.ReadAll(buf)
}

func TestGetNamespace(t *testing.T) {
	namespacePath := "default"

	reader := &ReadFiler{path: namespacePath} //default
	readFile = reader.ReadFile                //

	e := &Endpoint{}
	_ = e.GetNamespace(namespacePath)
	if namespacePath != e.Namespace {
		t.Errorf("wanted %s, got %s", namespacePath, e.Namespace)
	}
}
