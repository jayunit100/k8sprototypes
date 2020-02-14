package utils

import (
	"bytes"
//	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	v1net "k8s.io/api/networking/v1"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Kubernetes struct {
}

func (k *Kubernetes) Probe( ns string, pod  string, ns2  string, pod2  string, port int) int {
	ip := ""
	pod, err := k.Client().CoreV1().Pods(ns2).Get(pod2, metav1.GetOptions{})
	if err == nil {
		ip = pod.Status.PodIP
	} else {
		panic(err)
	}

	s:=fmt.Sprintf("%v:%v", ip, port)
	_,_,err = k.ExecWithOptions([]string{"curl", s, "-m", "5"}, ns, pod, "container")
	if err == nil {
		return 5
	}
	return -1
}

// ExecWithOptions executes a command in the specified container,
// returning stdout, stderr and error. `options` allowed for
// additional parameters to be passed.
func (k *Kubernetes) ExecWithOptions(Command []string, ns string, pod string, cname string) (string, string, error) {
	fmt.Println("Exec...")

	config, err := LoadConfig()
	Expect(err).NotTo(HaveOccurred(), "failed to load restclient config")

	req := f.ClientSet.Core().RESTClient().Post().
		Resource("pods").
		Name(pod).
		Namespace(ns).
		SubResource("exec").
		Param("container", cname)

	req.VersionedParams(&v1.PodExecOptions{
		Container: options.ContainerName,
		Command:   options.Command,
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}, api.ParameterCodec)

	var stdout, stderr bytes.Buffer
	err = execute("POST", req.URL(), config, options.Stdin, &stdout, &stderr, tty)
	return strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()), err
}

func (k *Kubernetes) Client() *kubernetes.Clientset {
	// TODO borrowed from e2e upstream.
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

	return clientset
}

func (k *Kubernetes) CreateNamespace(n string, labels map[string]string) (*v1.Namespace, error) {
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: n,
			Labels: labels,
		},
	}
	nsr, err := k.Client().CoreV1().Namespaces().Create(ns)
	return nsr, err
}

func (k *Kubernetes) CreateDeployment(ns, deploymentName string, replicas int32, labels map[string]string, image string) (*appsv1.Deployment, error) {
	zero := int64(0)
	d := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Labels:    labels,
			Namespace: ns,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: v1.PodSpec{
					TerminationGracePeriodSeconds: &zero,
					Containers: []v1.Container{
						{
							Name:            image,
							Image:           image,
							SecurityContext: &v1.SecurityContext{},
						},
					},
				},
			},
		},
	}

	return k.Client().AppsV1().Deployments("").Create(d)
}

func (k *Kubernetes) CreateNetworkPolicy(netpol *v1net.NetworkPolicy) (*v1net.NetworkPolicy,error) {
	np, err := k.Client().NetworkingV1().NetworkPolicies("").Create( netpol)
	return np, err
}
