package endpoint

import (
	"io/ioutil"

	apiCoreV1 "k8s.io/api/core/v1"
)

var readFile = ioutil.ReadFile

type Endpoint struct {
	Namespace string
	Svc       string
	Ips       []string
}

//GetNamespace returns current namespace
func (e *Endpoint) GetNamespace(namespacePath string) error {
	namespace, err := readFile(namespacePath)
	if err != nil {
		return err
	}
	e.Namespace = string(namespace)
	return nil
}

//GetAddresses returns ipaddrs for a service in a namespace
func (e *Endpoint) GetAddresses(endpoints *apiCoreV1.Endpoints) {
	ipaddrs := []string{}
	for _, subsets := range endpoints.Subsets {
		for _, addresses := range subsets.Addresses {
			ipaddrs = append(ipaddrs, addresses.IP)
		}
	}
	e.Ips = ipaddrs
}
