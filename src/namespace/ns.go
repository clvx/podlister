package namespace

import (
	"errors"

	apiCoreV1 "k8s.io/api/core/v1"
)

var (
	ErrorNotFoundNamespace = errors.New("ErrorNamespaceNotFound: Namespace does not exist in cluster")
)

type Namespace struct {
	Items []string
}

func iterateNamespaces(namespaces []apiCoreV1.Namespace) []string {
	var items []string
	for _, ns := range namespaces {
		if ns.ObjectMeta.Name == "kube-system" {
			continue
		}
		items = append(items, ns.ObjectMeta.Name)
	}
	return items
}

func namespaceLookup(ns string, items []string) bool {
	for _, item := range items {
		if item == ns {
			return true
		}
	}
	return false
}


//getNamespaces get all namespaces besides kube-system
func (n *Namespace) GetNamespaces(namespaces []apiCoreV1.Namespace) {
	/*
		//namespaces, err := cs.CoreV1().Namespaces().List(v1.ListOptions{})
		if err != nil {
			return err
		}
		//namespaces.Items //[]Namespace
	*/

	n.Items = iterateNamespaces(namespaces)
}

//getFilteredNamespace obtains a set of namespaces from all namespaces. It fails
//if a namespace is not found.
func (n *Namespace) GetFilteredNamespaces(namespaces []apiCoreV1.Namespace, nsList []string) error {

	var items []string
	/*
		namespaces, err := cs.CoreV1().Namespaces().List(v1.ListOptions{})
		if err != nil {
			return err
		}
	*/
	items = iterateNamespaces(namespaces)
	for _, ns := range nsList {
		if !namespaceLookup(ns, items) {
			return ErrorNotFoundNamespace
		}
	}
	return nil
}
