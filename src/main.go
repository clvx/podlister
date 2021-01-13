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
	end.GetAddresses(endpoints)

	//Rendering template
	var tpl bytes.Buffer
	tmpl := template.Must(template.ParseFiles(cfg.Template.Name))
	if err = tmpl.Execute(&tpl, end); err != nil {
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
