package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	networkv1 "k8s.io/api/networking/v1"
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
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	fmt.Println("Creating Ingress.......")
	pathType := "Exact"
	ing := &networkv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "prod-ingress",
			Namespace: "default",
		},
		Spec: networkv1.IngressSpec{
			DefaultBackend: &networkv1.IngressBackend{
				Service: &networkv1.IngressServiceBackend{
					Name: "vault",
					Port: networkv1.ServiceBackendPort{
						Name: "http",
					},
				},
			},
			Rules: []networkv1.IngressRule{
				{
					IngressRuleValue: networkv1.IngressRuleValue{
						HTTP: &networkv1.HTTPIngressRuleValue{
							Paths: []networkv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: (*networkv1.PathType)(&pathType),
									Backend: networkv1.IngressBackend{
										Service: &networkv1.IngressServiceBackend{
											Name: "vault-ui",
											Port: networkv1.ServiceBackendPort{
												Name: "http",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	result, err := clientset.NetworkingV1().Ingresses("default").Create(context.TODO(), ing, metav1.CreateOptions{})
	if err != nil {
		//handle error
		panic(err.Error())
	}
	fmt.Printf("Created Ingress %q.\n", result.GetObjectMeta().GetName())
}
