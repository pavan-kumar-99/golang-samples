package main

import (
	"context"
	"flag"
	"path/filepath"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
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
	secretname := "test-secret"
	secretClient := clientset.CoreV1().Secrets(namespace)
	logrus.Info("Checking for the existance of the secret")
	_, err = secretClient.Get(context.TODO(), secretname, metav1.GetOptions{})
	if err != nil {
		logrus.Warn("Secret not found. Creating it now")
	}
	logrus.Info("Secret found, Updating it to the latest values")

	//Update Logic with retry here
	//TODO, Add retry logic here and also add Update function

	err = secretClient.Delete(context.TODO(), secretname, metav1.DeleteOptions{})
	if err != nil {
		logrus.Error("Unable to delete the secret")
	}

	//Retry Ends here

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretname,
			Namespace: namespace,
		},
		StringData: map[string]string{
			"username": "Pavan",
			"Password": "Password123@",
		},
		Type: corev1.SecretType("Opqaue"),
	}
	secretc, err := secretClient.Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		logrus.Error("Failed to create a secret", err)
	}
	logrus.Info("Created secret ", secretc.ObjectMeta.Name)
}
