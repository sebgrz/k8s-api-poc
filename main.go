package main

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/util/homedir"
)

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		if err == rest.ErrNotInCluster {
			config, err = clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
			fmt.Println("Config is take from local machine")
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
		fmt.Println()
	}

	fmt.Println("Exec pod")
	pocPods := client.CoreV1().Pods("poc")
	pocPodsList, _ := pocPods.List(context.Background(), v1.ListOptions{})
	pod := pocPodsList.Items[0]

	req := client.CoreV1().RESTClient().Post().Resource("pods").Name(pod.Name).Namespace(pod.Namespace).SubResource("exec")
	fmt.Printf("Req to pod: %s\n", pod.Name)
	options := &corev1.PodExecOptions{
		Command: []string{"cat", "/etc/ssl/certs/ca-bundle.crt"},
		Stdout:  true,
	}
	req.VersionedParams(options, scheme.ParameterCodec)
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer

	err = exec.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdout: &buf,
	})

	if err != nil {
		panic(err)
	}

	// Stream to file in order to get full content
	fmt.Println(buf.String())
}
