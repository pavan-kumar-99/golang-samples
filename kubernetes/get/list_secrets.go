package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	var k kubernetes.Interface
	k, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	secret_list, err := k.CoreV1().Secrets("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, s := range secret_list.Items {
		ns := s.Namespace
		secret, err := k.CoreV1().Secrets(ns).Get(context.TODO(), s.GetName(), metav1.GetOptions{})
		if err != nil {
			panic(err)
		}
		fmt.Println(secret.Name)
	}
	fmt.Println(secret_list.Items[0].ObjectMeta.Namespace)
	fmt.Println(secret_list.Items[0].GetName())
	secret, err := k.CoreV1().Secrets("default").Get(context.TODO(), "vault-token-g729w", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	data := secret.Data
	for k, v := range data {
		val, err := json.Marshal(map[string]string{k: string(v)})
		if err != nil {
			panic(err)
		}
		fmt.Println(string(val))
	}

}
