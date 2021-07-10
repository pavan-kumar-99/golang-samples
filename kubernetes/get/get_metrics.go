package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"path/filepath"

	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	//	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
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
	//clientsetm, err := metricsv.NewForConfig(config)
	//if err != nil {
	//	panic(err)
	//}
	//podMetricsList, err := clientsetm.MetricsV1beta1().PodMetricses("metallb-system").List(context.TODO(), metav1.ListOptions{})
	//if err != nil {
	//panic(err)
	//}
	pods, err := clientset.CoreV1().Pods("metallb-system").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	//cpuUsagePercent := getCPUUsagePercent(podMetricsList.Items, pods.Items)
	jsonByte, err := json.Marshal(pods)
	if err != nil {
		//handle error
	}
	//fmt.Println(pods)
	fmt.Println(string(jsonByte))
	//fmt.Println("Pods", pods)
	//fmt.Println("Pods Items", pods.Items)
	//fmt.Println(podMetricsList.Items)
	//fmt.Println("CPU Percent is \n", cpuUsagePercent)
	//fmt.Println("Creating Pod.......")
	pod := &apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-pod",
		},
		Spec: apiv1.PodSpec{
			Containers: []apiv1.Container{
				{
					Name:  "nginx",
					Image: "nginx",
				},
			},
		},
	}
	result, err := clientset.CoreV1().Pods("default").Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		//handle error
		//panic(err.Error())
		fmt.Println("Pod alreasy exists")
	}
	fmt.Printf("Created Pod %q.\n", result.GetObjectMeta().GetName())
}

func getCPUUsagePercent(metrics []v1beta1.PodMetrics, pods []v1.Pod) []int32 {
	podResources := make(map[string]map[string]v1.ResourceRequirements, len(pods))
	for _, pod := range pods {
		containerMap := make(map[string]v1.ResourceRequirements, len(pod.Spec.Containers))
		for _, container := range pod.Spec.Containers {
			containerMap[container.Name] = container.Resources
		}
		podResources[pod.Name] = containerMap
		jsonByte, err := json.Marshal(podResources)
		if err != nil {
			//handle error
		}
		fmt.Println(string(jsonByte))
		jsonByte1, err := json.Marshal(containerMap)
		if err != nil {
			//handle error
		}
		fmt.Println(string(jsonByte1))
	}

	// calculate the max usage/request value for each pod by looking at all
	// containers in every pod.
	cpuUsagePercent := []int32{}
	for _, podMetrics := range metrics {
		if containerResources, ok := podResources[podMetrics.Name]; ok {
			podMax := int32(0)
			for _, container := range podMetrics.Containers {
				if resources, ok := containerResources[container.Name]; ok {
					requestedCPU := resources.Requests.Cpu().MilliValue()
					usageCPU := container.Usage.Cpu().MilliValue()
					if requestedCPU > 0 {
						usagePcnt := int32(100 * usageCPU / requestedCPU)
						if usagePcnt > podMax {
							podMax = usagePcnt
						}
					}
				}
			}
			cpuUsagePercent = append(cpuUsagePercent, podMax)
		}
	}
	return cpuUsagePercent
}

