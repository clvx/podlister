package main

import (
	"bytes"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Endpoint struct {
	Svc string
	Ips []string
}

func getAddresses(cs *kubernetes.Clientset, endp *Endpoint) error {
	ipaddrs := []string{}
	endpoints, err := cs.CoreV1().Endpoints("pd").Get(endp.Svc, v1.GetOptions{})
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
	endpoint := &Endpoint{Svc: "pd-swarm"}
	err = getAddresses(clientset, endpoint)
	if err != nil {
		log.Println(err)
	}

	//Rendering template
	var tpl bytes.Buffer
	tmpl := template.Must(template.ParseFiles("index.template"))
	if err = tmpl.Execute(&tpl, endpoint); err != nil {
		log.Println(err)
	}

	// Write function to upload to object storage
	key := os.Getenv("SPACES_KEY")
	secret := os.Getenv("SPACES_SECRET")

	s3config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:    aws.String("https://nyc3.digitaloceanspaces.com"),
		Region:      aws.String("us-east-1"),
	}
	newSession := session.New(s3config)
	s3Client := s3.New(newSession)

	//TODO: manage inputs
	object := s3.PutObjectInput{
		Bucket: aws.String("pd-swarm"),
		Key:    aws.String("index.html"),
		Body:   strings.NewReader(tpl.String()),
		ACL:    aws.String("public-read"),
	}
	_, err = s3Client.PutObject(&object)
	if err != nil {
		log.Println(err.Error())
	}

}
