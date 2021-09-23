package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
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
	fmt.Println("Creating Service.......")
	svc := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "prod-svc",
			Namespace: "default",
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				{
					Name: "http",
					Port: int32(80),
					TargetPort: intstr.IntOrString{
						IntVal: 80,
					},
				},
			},
			Selector: map[string]string{
				"app": "nginx-deploy",
			},
			Type: "ClusterIP",
		},
	}
	result, err := clientset.CoreV1().Services("default").Create(context.TODO(), svc, metav1.CreateOptions{})
	if err != nil {
		//handle error
		panic(err.Error())
	}
	fmt.Printf("Created Service %q.\n", result.GetObjectMeta().GetName())
}
