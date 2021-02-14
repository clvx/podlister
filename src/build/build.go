package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

//Obtain credentials to connect to the cluster

type cluster struct {
	clientSet *kubernetes.Clientset
}

//createNamespaces creates a number of namespaces
func (c *cluster) createNamespaces(namespaces ...string) (err error) {
	for _, n := range namespaces {
		item := &v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: n,
			},
		}

		ns, err := c.clientSet.CoreV1().Namespaces().Create(item)
		if err != nil {
			return err
		}
		fmt.Println(ns.Status)
	}
	return nil
}

func (c *cluster) deleteNamespaces(namespaces ...string) (err error) {
	for _, n := range namespaces {
		err := c.clientSet.CoreV1().Namespaces().Delete(n, &metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *cluster) createPods(pods ...string) (err error) {
	for _, p := range pods {
		item := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      p,
				Namespace: "default",
				Labels: map[string]string{
					"app": "demo",
				},
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:            p,
						Image:           "busybox",
						ImagePullPolicy: v1.PullIfNotPresent,
						Command: []string{
							"sleep",
							"3600",
						},
					},
				},
			},
		}
		pod, err := c.clientSet.CoreV1().Pods("default").Create(item)
		if err != nil {
			return err
		}
		fmt.Println(pod.Status)
	}
	return nil
}

func (c *cluster) deletePods(pods ...string) (err error) {
	for _, p := range pods {
		err := c.clientSet.CoreV1().Pods("default").Delete(p, &metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

//createServices creates a regular service with clusterIP as default
func (c *cluster) createServices(svcs ...string) (err error) {
	for _, s := range svcs {
		item := &v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      s,
				Namespace: "default",
				Labels: map[string]string{
					"app": "demo",
				},
			},
			Spec: v1.ServiceSpec{
				Ports: []v1.ServicePort{
					{
						Name:     "http",
						Protocol: "TCP",
						Port:     8080,
						TargetPort: intstr.IntOrString{
							IntVal: 8080,
						},
					},
				},
				Selector: map[string]string{
					"app": "demo",
				},
			},
		}
		svc, err := c.clientSet.CoreV1().Services("default").Create(item)
		if err != nil {
			return err
		}
		fmt.Println(svc.Status)
	}
	return nil
}

func (c *cluster) deleteServices(svcs ...string) (err error)            {
	for _, s := range svcs {
		err := c.clientSet.CoreV1().Services("default").Delete(s, &metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}
	return nil

}
func (c *cluster) createIngress()             {}
func (c *cluster) createIstioVirtualService() {}
func (c *cluster) createIstioIngressGateway() {}
func (c *cluster) createIstioEgressGateway()  {}

//installIstio installs Istio
func (c *cluster) installIstio() {}

//installIngress installs an ingress according to the ingress type
func (c *cluster) installIngress() {}

func buildCluster(c cluster) {
	//createNamespaces
	err := c.createNamespaces([]string{"foo", "bar"}...) //passing variadic arguments
	if err != nil {
		log.Fatalf("error %v", err)
	}

	//Create Pod
	err = c.createPods([]string{"fooz", "barz"}...)
	if err != nil {
		log.Fatalf("error %v", err)
	}
	//Create Service
	err = c.createServices([]string{"fooz", "barz"}...)
	if err != nil {
		log.Fatalf("error %v", err)
	}
}

func cleanCluster(c cluster) {
	err := c.deleteNamespaces([]string{"foo", "bar"}...) //passing variadic arguments
	if err != nil {
		log.Fatalf("error %v", err)
	}

	err = c.deletePods([]string{"fooz", "barz"}...)
	if err != nil {
		log.Fatalf("error %v", err)
	}

	err = c.deletePods([]string{"fooz", "barz"}...)
	if err != nil {
		log.Fatalf("error %v", err)
	}
}

func main() {
	//Obtaining kubeconfig
	kubeconfig := filepath.Join("~", ".kube", "config")
	if envvar := os.Getenv("KUBECONFIG"); len(envvar) > 0 {
		kubeconfig = envvar
	}
	kconfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Printf("Kubeconfig cannot be loaded: %v\n", err)
		os.Exit(2)
	}
	// Creates clientset
	clientset, err := kubernetes.NewForConfig(kconfig)
	if err != nil {
		panic(err.Error())
	}
	c := cluster{clientset}

	//buildCluster(c)
	cleanCluster(c)
}
