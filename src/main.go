package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	conf "podlister/config"

	"podlister/namespace"

	"github.com/spf13/viper"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	namespacePath = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
	configPath    = "./config/config.yaml"
)

func loadConfig(path string) (config conf.Config, err error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func main() {
	//Loading configs
	cfg, err := loadConfig("./config")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	//Obtaining kubeconfig
	kconfig, err := rest.InClusterConfig()
	if err != nil {
		kubeconfig := filepath.Join("~", ".kube", "config")
		if envvar := os.Getenv("KUBECONFIG"); len(envvar) > 0 {
			kubeconfig = envvar
		}
		kconfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Printf("Kubeconfig cannot be loaded: %v\n", err)
			os.Exit(2)
		}
	}
	// Creates clientset
	clientset, err := kubernetes.NewForConfig(kconfig)
	if err != nil {
		panic(err.Error())
	}
	/*
	   Inputs:
	   ns{*; *}
	   ns{*; label}
	   ns{ns1; *}
	   ns{ns; label}
	*/

	ns := namespace.Namespace{}
	if len(cfg.Namespace) > 0 {				//filter given namespaces
		nsClusterFiltered, err := clientset.CoreV1().Namespaces().List(v1.ListOptions{})
		if err != nil {
			log.Fatalf("Error %v", err)
		}
		ns.GetFilteredNamespaces(nsClusterFiltered.Items, cfg.Namespace)
	} else {								//get all namespaces
		nsClusterList, err := clientset.CoreV1().Namespaces().List(v1.ListOptions{})
		if err != nil {
			log.Fatalf("Error %v", err)
		}
		ns.GetNamespaces(nsClusterList.Items)
	}

}
