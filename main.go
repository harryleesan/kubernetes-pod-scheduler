package main

import (
	"errors"
	"fmt"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
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

		if os.Getenv("SCALE") == "scaleUp" {
			fmt.Println("Scaling up...")
		} else if os.Getenv("SCALE") == "scaleDown" {
			fmt.Println("Scaling down...")
		} else {
			return errors.New("Incorrect env var for SCALE!")
		}

		i, intConvError := strconv.ParseInt(d.ObjectMeta.Annotations[os.Getenv("SCALE")], 10, 32)
		if intConvError != nil {
			return intConvError
		}
		_, err := clientset.ExtensionsV1beta1().Deployments(namespace).UpdateScale(name, &v1beta1.Scale{
			TypeMeta:   d.TypeMeta,
			ObjectMeta: d.ObjectMeta,
			Spec:       v1beta1.ScaleSpec{int32(i)},
		})
		if err != nil {
			return err
		}
		fmt.Printf("Scaled pods to: %v\n", int32(i))
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

	// use the current context in kubeconfig
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}
