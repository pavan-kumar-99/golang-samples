package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	resource "k8s.io/apimachinery/pkg/api/resource"
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
	var replicas int32 = 2
	stsspec := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "golang-sts",
			Namespace: "default",
			Labels: map[string]string{
				"app":        "go-sts",
				"managed-by": "go-sts",
			},
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":   "golang-client",
					"owned": "golang-client",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: "golang-pod",
					Labels: map[string]string{
						"app":   "golang-client",
						"owned": "golang-client",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "main-container",
							Image: "nginx",
						},
						{
							Name:    "sidecar-container",
							Image:   "busybox",
							Command: []string{"sleep", "10000"},
						},
					},
				},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "golang-client-volume",
					},
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes: []corev1.PersistentVolumeAccessMode{
							corev1.ReadWriteOnce,
						},
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceName(corev1.ResourceStorage): resource.MustParse("1Gi"),
							},
						},
					},
				},
			},
		},
	}
	sts, err := k.AppsV1().StatefulSets("default").Create(context.TODO(), stsspec, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(sts.GetName())
}
