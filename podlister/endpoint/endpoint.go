package endpoint

import (
	"io/ioutil"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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
// TODO:Add GetAddresses test
func (e *Endpoint) GetAddresses(cs *kubernetes.Clientset) error {
	ipaddrs := []string{}
	endpoints, err := cs.CoreV1().Endpoints(e.Namespace).Get(e.Svc, v1.GetOptions{})
	if err != nil {
		return err
	}
	for _, subsets := range endpoints.Subsets {
		for _, addresses := range subsets.Addresses {
			ipaddrs = append(ipaddrs, addresses.IP)
		}
	}
	e.Ips = ipaddrs
	return nil
}
