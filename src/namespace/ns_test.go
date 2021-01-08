package namespace

import (
	"reflect"
	"testing"

	apiCoreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSetNamespace(t *testing.T) {
	namespace := "default"

	n := Namespace{}
	n.setNamespace(namespace)
	if namespace != n.Items[0] {
		t.Errorf("wanted %s, got %s", namespace, n.Items[0])
	}
}

func TestGetNamespaces(t *testing.T) {
	n := Namespace{}
	namespaces := []apiCoreV1.Namespace{
		{
			metav1.TypeMeta{},
			metav1.ObjectMeta{
				Name: "foo",
			},
			apiCoreV1.NamespaceSpec{},
			apiCoreV1.NamespaceStatus{},
		},
		{
			metav1.TypeMeta{},
			metav1.ObjectMeta{
				Name: "bar",
			},
			apiCoreV1.NamespaceSpec{},
			apiCoreV1.NamespaceStatus{},
		},
		{
			metav1.TypeMeta{},
			metav1.ObjectMeta{
				Name: "kube-system",
			},
			apiCoreV1.NamespaceSpec{},
			apiCoreV1.NamespaceStatus{},
		},
	}
	t.Run("All namespaces", func(t *testing.T) {

		expected := Namespace{[]string{"foo", "bar"}}

		n.getNamespaces(namespaces)
		if !reflect.DeepEqual(expected.Items, n.Items) {
			t.Errorf("wanted %v, got %v", expected.Items, n.Items)
		}
	})
	t.Run("Namespaces with kube-system", func(t *testing.T) {
		expected := Namespace{[]string{"foo", "bar"}}

		n.getNamespaces(namespaces)

		if !reflect.DeepEqual(expected.Items, n.Items) {
			t.Errorf("wanted %v, got %v", expected, n.Items)
		}

	})
}

func TestFilteredNamespaces(t *testing.T) {
	n := Namespace{}
	namespaces := []apiCoreV1.Namespace{
		{
			metav1.TypeMeta{},
			metav1.ObjectMeta{
				Name: "foo",
			},
			apiCoreV1.NamespaceSpec{},
			apiCoreV1.NamespaceStatus{},
		},
		{
			metav1.TypeMeta{},
			metav1.ObjectMeta{
				Name: "bar",
			},
			apiCoreV1.NamespaceSpec{},
			apiCoreV1.NamespaceStatus{},
		},
		{
			metav1.TypeMeta{},
			metav1.ObjectMeta{
				Name: "kube-system",
			},
			apiCoreV1.NamespaceSpec{},
			apiCoreV1.NamespaceStatus{},
		},
	}
	t.Run("filtered namespace", func(t *testing.T){
		nsList := []string{"foo"}
		err := n.getFilteredNamespaces(namespaces, nsList)
		if err != nil {
			t.Errorf("wanted %v, got %v", nil, err)
		}
	})
	t.Run("Non existent namespace", func(t *testing.T){
		nsList := []string{"foobar"}
		err := n.getFilteredNamespaces(namespaces, nsList)
		if err != ErrorNotFoundNamespace {
			t.Errorf("wanted %s, got %v", ErrorNotFoundNamespace, err)
		}
	})

}
