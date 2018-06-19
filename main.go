package main

import (
	"flag"
	"fmt"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"strconv"
)

var clientset *kubernetes.Clientset

func scale(d v1beta1.Deployment) error {
	if d.ObjectMeta.Annotations["scaleDown"] != "" && d.ObjectMeta.Annotations["scaleUp"] != "" {
		fmt.Printf("Deployment: %#v\nNamespace: %#v\n", d.ObjectMeta.Name, d.ObjectMeta.Namespace)
		fmt.Printf("Scale Down: %#v\n", d.ObjectMeta.Annotations["scaleDown"])
		fmt.Printf("Scale Up: %#v\n", d.ObjectMeta.Annotations["scaleUp"])
		namespace := d.ObjectMeta.Namespace
		name := d.ObjectMeta.Name
		i, intConvError := strconv.ParseInt(d.ObjectMeta.Annotations[os.Getenv("SCALE")], 10, 32)
		if intConvError != nil {
			return intConvError
		}
		scaleRes, err := clientset.ExtensionsV1beta1().Deployments(namespace).UpdateScale(name, &v1beta1.Scale{
			TypeMeta:   d.TypeMeta,
			ObjectMeta: d.ObjectMeta,
			Spec:       v1beta1.ScaleSpec{int32(i)},
		})
		if err != nil {
			return err
		}
		fmt.Printf("Scaling pods to: %v\n", scaleRes.Status.Replicas)
		fmt.Println("")
	}
	return nil
}

func main() {
	kubeClientSetUp()
	deploymentlist, err := clientset.ExtensionsV1beta1().Deployments("").List(metav1.ListOptions{})
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Printf("#%v", deploymentlist)
	for _, deployment := range deploymentlist.Items {
		scaleError := scale(deployment)
		if scaleError != nil {
			fmt.Printf("Error! %v", scaleError)
		}
	}

}

// kubeClientSetUp() creates the client used to connect to the kubernetes
// cluster.
func kubeClientSetUp() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	fmt.Printf("Kubeconfig used: %s\n\n", *kubeconfig)
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
