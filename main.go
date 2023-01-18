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

	namespaces := client.CoreV1().Namespaces()
	namespaceList, err := namespaces.List(context.Background(), v1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, namespace := range namespaceList.Items {
		pods := client.CoreV1().Pods(namespace.Name)
		podList, err := pods.List(context.Background(), v1.ListOptions{})
		if err != nil {
			panic(err)
		}

		fmt.Printf("Pods of namespace %s:\n", namespace.Name)
		for _, i := range podList.Items {
			fmt.Printf("Name: %s Namespace: %s\n", i.Name, i.Namespace)
		}
		println()
	}
}
