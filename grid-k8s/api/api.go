package api

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	//"time"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	k8errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	watch "k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

const (
	K8State_All        string = ""
	K8State_Waiting           = "Waiting"
	K8State_Running           = "Running"
	K8State_Terminated        = "Terminated"
)

type KubeApi struct {
	clientset *kubernetes.Clientset
}

func NewKubeApi() *KubeApi {
	k := new(KubeApi)
	return k
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func (k *KubeApi) InitKube() error {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return err
	}
	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	k.clientset = clientset
	return nil
}

func (k *KubeApi) GetNodeList() (*v1.NodeList, error) {
	nodeList, err := k.clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	return nodeList, err
}

func (k *KubeApi) GetNodes(tag string) ([]*v1.Node, error) {
	nodeList, err := k.GetNodeList()
	if err != nil {
		return nil, err
	}
	var nodes []*v1.Node = nil
	if tag == "" {
		nodes = make([]*v1.Node, 0, len(nodeList.Items))
	} else {
		nodes = make([]*v1.Node, 0, 0)
	}
	num := len(nodeList.Items)
	for i := 0; i < num; i++ {
		node := nodeList.Items[i]
		if tag == "" || strings.HasPrefix(node.Name, tag) {
			nodes = append(nodes, &node)
		}
	}
	return nodes, nil
}

func (k *KubeApi) GetNode(name string) (*v1.Node, error) {
	node, err := k.clientset.CoreV1().Nodes().Get(context.TODO(), name, metav1.GetOptions{})
	if k8errors.IsNotFound(err) {
		return nil, errors.New(fmt.Sprintf("Node %s not found\n", name))
	} else if statusError, isStatus := err.(*k8errors.StatusError); isStatus {
		return nil, errors.New(fmt.Sprintf("Error getting node %s : %v\n",
			name, statusError.ErrStatus.Message))
	} else if err != nil {
		return nil, err
	}
	return node, err
}

func (k *KubeApi) GetPodList(namespace string) (*v1.PodList, error) {
	podList, err := k.clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	return podList, err
}

func (k *KubeApi) GetPods(namespace string, tag string) ([]*v1.Pod, error) {
	podList, err := k.GetPodList(namespace)
	if err != nil {
		return nil, err
	}
	var pods []*v1.Pod = nil
	if tag == "" {
		pods = make([]*v1.Pod, 0, len(podList.Items))
	} else {
		pods = make([]*v1.Pod, 0, 0)
	}
	num := len(podList.Items)
	for i := 0; i < num; i++ {
		pod := podList.Items[i]
		if tag == "" || strings.HasPrefix(pod.Name, tag) {
			pods = append(pods, &pod)
		}
	}
	return pods, nil
}

func (k *KubeApi) GetPod(namespace string, name string) (*v1.Pod, error) {
	pod, err := k.clientset.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if k8errors.IsNotFound(err) {
		return nil, errors.New(fmt.Sprintf("Pod %s in namespace %s not found\n", name, namespace))
	} else if statusError, isStatus := err.(*k8errors.StatusError); isStatus {
		return nil, errors.New(fmt.Sprintf("Error getting pod %s in namespace %s: %v\n",
			name, namespace, statusError.ErrStatus.Message))
	} else if err != nil {
		return nil, err
	}
	return pod, err
}

func (k *KubeApi) GetServiceList(namespace string) (*v1.ServiceList, error) {
	serviceList, err := k.clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	return serviceList, err
}

func (k *KubeApi) GetServices(namespace string, tag string) ([]*v1.Service, error) {
	serviceList, err := k.GetServiceList(namespace)
	if err != nil {
		return nil, err
	}
	var services []*v1.Service = nil
	if tag == "" {
		services = make([]*v1.Service, 0, len(serviceList.Items))
	} else {
		services = make([]*v1.Service, 0, 0)
	}
	num := len(serviceList.Items)
	for i := 0; i < num; i++ {
		service := serviceList.Items[i]
		if tag == "" || strings.HasPrefix(service.Name, tag) {
			services = append(services, &service)
		}
	}
	return services, nil
}

func (k *KubeApi) GetService(namespace string, name string) (*v1.Service, error) {
	service, err := k.clientset.CoreV1().Services(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if k8errors.IsNotFound(err) {
		return nil, errors.New(fmt.Sprintf("Service %s in namespace %s not found\n", name, namespace))
	} else if statusError, isStatus := err.(*k8errors.StatusError); isStatus {
		return nil, errors.New(fmt.Sprintf("Error getting service %s in namespace %s: %v\n",
			name, namespace, statusError.ErrStatus.Message))
	} else if err != nil {
		return nil, err
	}
	return service, err
}

func (k *KubeApi) GetDeploymentList(namespace string) (*appsv1.DeploymentList, error) {
	deploymentList, err := k.clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	return deploymentList, err
}

func (k *KubeApi) GetDeployments(namespace string, tag string) ([]*appsv1.Deployment, error) {
	deploymentList, err := k.GetDeploymentList(namespace)
	if err != nil {
		return nil, err
	}
	var deployments []*appsv1.Deployment = nil
	if tag == "" {
		deployments = make([]*appsv1.Deployment, 0, len(deploymentList.Items))
	} else {
		deployments = make([]*appsv1.Deployment, 0, 0)
	}
	num := len(deploymentList.Items)
	for i := 0; i < num; i++ {
		deployment := deploymentList.Items[i]
		if tag == "" || strings.HasPrefix(deployment.Name, tag) {
			deployments = append(deployments, &deployment)
		}
	}
	return deployments, nil
}

func (k *KubeApi) GetDeployment(namespace string, name string) (*appsv1.Deployment, error) {
	deployment, err := k.clientset.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if k8errors.IsNotFound(err) {
		return nil, errors.New(fmt.Sprintf("Service %s in namespace %s not found\n", name, namespace))
	} else if statusError, isStatus := err.(*k8errors.StatusError); isStatus {
		return nil, errors.New(fmt.Sprintf("Error getting deployment %s in namespace %s: %v\n",
			name, namespace, statusError.ErrStatus.Message))
	} else if err != nil {
		return nil, err
	}
	return deployment, err
}

func (k *KubeApi) GetDaemonSetList(namespace string) (*appsv1.DaemonSetList, error) {
	daemonSetList, err := k.clientset.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
	return daemonSetList, err
}

func (k *KubeApi) GetDaemonSets(namespace string, tag string) ([]*appsv1.DaemonSet, error) {
	daemonSetList, err := k.GetDaemonSetList(namespace)
	if err != nil {
		return nil, err
	}
	var daemonSets []*appsv1.DaemonSet = nil
	if tag == "" {
		daemonSets = make([]*appsv1.DaemonSet, 0, len(daemonSetList.Items))
	} else {
		daemonSets = make([]*appsv1.DaemonSet, 0, 0)
	}
	num := len(daemonSetList.Items)
	for i := 0; i < num; i++ {
		daemonSet := daemonSetList.Items[i]
		if tag == "" || strings.HasPrefix(daemonSet.Name, tag) {
			daemonSets = append(daemonSets, &daemonSet)
		}
	}
	return daemonSets, nil
}

func (k *KubeApi) GetDaemonSet(namespace string, name string) (*appsv1.DaemonSet, error) {
	daemonSet, err := k.clientset.AppsV1().DaemonSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if k8errors.IsNotFound(err) {
		return nil, errors.New(fmt.Sprintf("Service %s in namespace %s not found\n", name, namespace))
	} else if statusError, isStatus := err.(*k8errors.StatusError); isStatus {
		return nil, errors.New(fmt.Sprintf("Error getting daemonSet %s in namespace %s: %v\n",
			name, namespace, statusError.ErrStatus.Message))
	} else if err != nil {
		return nil, err
	}
	return daemonSet, err
}

func (k *KubeApi) WatchNodeList() (watch.Interface, error) {
	watcher, err := k.clientset.CoreV1().Nodes().Watch(context.TODO(), metav1.ListOptions{})
	return watcher, err
}

func (k *KubeApi) WatchPodList(namespace string) (watch.Interface, error) {
	watcher, err := k.clientset.CoreV1().Pods(namespace).Watch(context.TODO(), metav1.ListOptions{})
	return watcher, err
}

func (k *KubeApi) WatchServiceList(namespace string) (watch.Interface, error) {
	watcher, err := k.clientset.CoreV1().Services(namespace).Watch(context.TODO(), metav1.ListOptions{})
	return watcher, err
}

func (k *KubeApi) WatchDeploymentList(namespace string) (watch.Interface, error) {
	watcher, err := k.clientset.AppsV1().Deployments(namespace).Watch(context.TODO(), metav1.ListOptions{})
	return watcher, err
}

func (k *KubeApi) WatchDaemonSetList(namespace string) (watch.Interface, error) {
	watcher, err := k.clientset.AppsV1().DaemonSets(namespace).Watch(context.TODO(), metav1.ListOptions{})
	return watcher, err
}
