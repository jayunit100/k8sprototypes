package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	fmt.Println("hi")
	//_ = tests.Scenario{}
	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	api := clientset.CoreV1()
	podList, err := api.Pods("").List(v1.ListOptions{})
	if err == nil {
		for _, p := range podList.Items {
			fmt.Println("%s", p)
		}
	} else {
		panic(err)
		fmt.Println("failed")
	}
	// Create a K8sScenario, m
	/**
			whitelist := map[string]bool{}...
			In m:
			 For namespaces (a b c)
				m.WithDeployments (a b c)

			p = Use the builder to make a network policy.
			In m:
				ApplyPolicies(p)

			// validate policies on whitelist
			r := newReachableMatrix()
			for namespaces:
				for pods:
					for namespaces:
						for pods:
							ReachableMatrix.add(n1, pod1, n2, pod2, m.Probe(n1, pod1, n2, pod2)

			testResult := false

	**/

}
