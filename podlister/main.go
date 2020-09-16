package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"

	"podlister/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ilyakaznacheev/cleanenv"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	namespacePath = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
	configPath    = "./config/config.yaml"
)


func getNamespace(cs *kubernetes.Clientset, endp *Endpoint) error {
	namespace, err := ioutil.ReadFile(namespacePath)
	if err != nil {
		return err 
	}
	endp.Namespace = string(namespace)
	return nil
}

func getAddresses(cs *kubernetes.Clientset, endp *Endpoint) error {
	ipaddrs := []string{}
	endpoints, err := cs.CoreV1().Endpoints(endp.Namespace).Get(endp.Svc, v1.GetOptions{})
	if err != nil {
		return err
	}
	for _, subsets := range endpoints.Subsets {
		for _, addresses := range subsets.Addresses {
			ipaddrs = append(ipaddrs, addresses.IP)
		}
	}
	endp.Ips = ipaddrs
	return nil
}

func main() {
	var cfg config.Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Println(err)
		os.Exit(2)
	}

	// Creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// Creates clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Obtaining endpoints
	endpoint := &Endpoint{Svc: cfg.Service.Name}
	err = getNamespace(clientset, endpoint)
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}
	err = getAddresses(clientset, endpoint)
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}

	//Rendering template
	var tpl bytes.Buffer
	tmpl := template.Must(template.ParseFiles(cfg.Template.Name))
	if err = tmpl.Execute(&tpl, endpoint); err != nil {
		log.Println(err)
		os.Exit(2)
	}

	//Uploading data
	s3config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(cfg.Bucket.Key, cfg.Bucket.Secret, ""),
		Endpoint:    aws.String(cfg.Bucket.URL),
		Region:      aws.String(cfg.Bucket.Region),
	}
	newSession := session.New(s3config)
	s3Client := s3.New(newSession)

	object := s3.PutObjectInput{
		Bucket: aws.String(cfg.Bucket.Name),
		Key:    aws.String(cfg.Template.Output),
		Body:   strings.NewReader(tpl.String()),
		ACL:    aws.String(cfg.Bucket.Privilege),
	}
	_, err = s3Client.PutObject(&object)
	if err != nil {
		log.Println(err.Error())
		os.Exit(2)
	} else {
		log.Printf("%s Uploaded successfully", cfg.Template.Output)
	}
}
