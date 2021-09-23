package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	autoscale "k8s.io/api/autoscaling/v1"
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
	fmt.Println("Creating HPA.......")
	var min int32 = 1
	var max int32 = 2
	var targetCPU int32 = 1
	hpa := &autoscale.HorizontalPodAutoscaler{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{Name: "prod-hpa", Namespace: "default"},
		Spec: autoscale.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscale.CrossVersionObjectReference{
				Kind:       "Deployment",
				Name:       "ingress-nginx-controller",
				APIVersion: "apps/v1",
			},
			MinReplicas:                    &min,
			MaxReplicas:                    max,
			TargetCPUUtilizationPercentage: &targetCPU,
		},
	}
	result, err := clientset.AutoscalingV1().HorizontalPodAutoscalers("default").Create(context.TODO(), hpa, metav1.CreateOptions{})
	if err != nil {
		//handle error
		panic(err.Error())
	}
	fmt.Printf("Created HPA %q.\n", result.GetObjectMeta().GetName())
}
