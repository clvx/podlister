package endpoint

import (
	"bytes"
	"io/ioutil"
	"testing"

	apiCoreV1 "k8s.io/api/core/v1"
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

	reader := &ReadFiler{path: namespacePath}
	readFile = reader.ReadFile

	e := &Endpoint{}
	_ = e.GetNamespace(namespacePath)
	if namespacePath != e.Namespace {
		t.Errorf("wanted %s, got %s", namespacePath, e.Namespace)
	}
}

func TestGetAddresses(t *testing.T) {

	endpoints := &apiCoreV1.Endpoints{
		Subsets: []apiCoreV1.EndpointSubset{
			{
				Addresses: []apiCoreV1.EndpointAddress{
					{IP: "10.0.0.1"},
					{IP: "10.0.0.2"},
				},
			},
		},
	}

	e := &Endpoint{}
	e.GetAddresses(endpoints)
	addrs := endpoints.Subsets[0].Addresses
	if len(e.Ips) == len(addrs) {
		for i := range addrs {
			if e.Ips[i] != addrs[i].IP {
				t.Errorf("wanted %s, got %s", addrs[i].IP, e.Ips[i])
			}
		}
	} else {
		t.Errorf("wanted %v, got %v", addrs, e.Ips)
	}
}
