package main

import (
	"context"
	"fmt"
	"path/filepath"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		if err == rest.ErrNotInCluster {
			config, err = clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
			println("Config is take from local machine")
		}
		if err != nil {
			panic(err)
		}
	}
	
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	pods := client.CoreV1().Pods("default")
	podList, err := pods.List(context.Background(), v1.ListOptions{})
	if err != nil {
		panic(err)
	}

	println("Pods:")
	for _, i := range podList.Items {
	fmt.Printf("Name: %s\n", i.Name)
	}
}
