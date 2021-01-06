package namespace 

import "errors"

var (
	ErrorNamespaceNotFound = errors.New("ErrorNamespaceNotFound: Namespace does not exist in cluster")
)

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	 apiCoreV1 "k8s.io/api/core/v1"
)

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

func namespaceLookup(ns string, items []string) bool{
	for _, item := range items {
		if item == ns {
			return true	
		}
	}
	return false
}

type Namespace struct {
	Items []string
}

//getNamespaces gets all namespaces besides kube-system
func (n *Namespace) getNamespaces(cs *kubernetes.Clientset) error {
	namespaces, err := cs.CoreV1().Namespaces().List(v1.ListOptions{})
	if err != nil {
		return err
	}
	n.Items = iterateNamespaces(namespaces.Items)
	return nil
}

//getNamespace obtains the desired namespace
func (n *Namespace) getNamespace(cs *kubernetes.Clientset, ns string) error {
	namespace, err := cs.CoreV1().Namespaces().Get(ns, v1.GetOptions{})
	if err != nil {
		return err
	}
	n.Items = append(n.Items, namespace.ObjectMeta.Name)
	return nil
}

//getFilteredNamespace obtains a set of namespaces from all namespaces. It fails
//if a namespace is not found.
func (n *Namespace) getFilteredNamespaces(cs *kubernetes.Clientset, nsList []string) error {

	var items []string
	namespaces, err := cs.CoreV1().Namespaces().List(v1.ListOptions{})
	if err != nil {
		return error.Error()
	}
	items = iterateNamespaces(namespaces.Items)
	for _, ns := nsList {
		if ! namespaceLookup(ns, items) {
			return ErrorNamespaceNotFound
		}
	}
	return nil
}
