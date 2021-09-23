package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
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
	namespace := "default"
	deploynaamae := "go-deployment"
	fmt.Println("Checking for a existing deployment")
	_, err = clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploynaamae, metav1.GetOptions{})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"namespace":  namespace,
			"Deployment": deploynaamae,
		}).Info("Deploymet not found in namespace ")
	}
	logrus.Info("Deleting Deployment before creating")
	err = clientset.AppsV1().Deployments(namespace).Delete(context.TODO(), deploynaamae, metav1.DeleteOptions{})
	if err != nil {
		logrus.Info("Errof failed deleting deploy", err.Error())
	}
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploynaamae,
			Namespace: namespace,
			Annotations: map[string]string{
				"kv-inject":  "true",
				"init-only":  "false",
				"consul-url": "https://consul-prod.pavan.com",
				"Port":       "9200",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":      "golang",
					"deployed": "go-cleint",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":      "golang",
						"deployed": "go-cleint",
					},
				},
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						{
							Name:    "echo-init-1",
							Image:   "busybox",
							Command: []string{"echo", "In Init Container"},
							Resources: corev1.ResourceRequirements{
								Limits: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("1"),
									"memory": resource.MustParse("1Gi"),
								},
								Requests: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("100m"),
									"memory": resource.MustParse("90Mi"),
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:    "main-container",
							Image:   "busybox",
							Command: []string{"sleep", "10000"},
							Resources: corev1.ResourceRequirements{
								Limits: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("1"),
									"memory": resource.MustParse("1Gi"),
								},
								Requests: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("100m"),
									"memory": resource.MustParse("90Mi"),
								},
							},
						},
					},
				},
			},
		},
	}
	create_deploy, err := clientset.AppsV1().Deployments("default").Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		logrus.Panic("Error creating deployment", err.Error())
	}
	annotations := create_deploy.ObjectMeta.Annotations
	logrus.Info(annotations)
	extract(annotations)
}

func int32Ptr(i int32) *int32 { return &i }

func extract(c map[string]string) {
	if consul_url, ok := c["consul-url"]; ok {
		fmt.Println("Consul url is", consul_url)
	}
}
