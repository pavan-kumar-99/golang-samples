package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"time"

	//"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"

	//appslisters "k8s.io/client-go/listers/apps/v1"

	"k8s.io/client-go/tools/cache"
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
	sif := informers.NewSharedInformerFactory(clientset, 0*time.Second)
	podinfo := sif.Core().V1().Pods()
	podinfo.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    adddeploy,
		DeleteFunc: deletedeploy,
	})

	sif.Start(wait.NeverStop)
	sif.WaitForCacheSync(wait.NeverStop)
	//time.Sleep(10 * time.Second)
	//deploy, _ := appslisters.DeploymentLister.Deployments("").Get("go-deployment")
	pods, err := podinfo.Lister().Pods("default").Get("vault-0")
	if err != nil {
		panic(err)
	}
	fmt.Println(pods.ObjectMeta.Name)
	//fmt.Println(deploy)
}

func adddeploy(o interface{}) {
	//fmt.Println(o)
	fmt.Println("In add function")
}

func deletedeploy(j interface{}) {
	//fmt.Println(j)
	fmt.Println("In Delete function")
}
