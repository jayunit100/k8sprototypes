package utils

import (
	"bytes"
	"io"
	"net/url"
	//	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	v1net "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

type Kubernetes struct {
}

func (k *Kubernetes) Probe(ns1 string, pod1 string, ns2 string, pod2 string, port int) int {
	//ip := ""
	pod, err := k.Client().CoreV1().Pods(ns2).Get(pod2, metav1.GetOptions{})
	//.CoreV1().Pods(ns2).Get(pod2, metav1.GetOptions{})
	if err == nil {
		//ip = pod.Status.PodIP
	} else {
		panic(err)
	}
	fmt.Println("Pod ip", pod)
	//s := fmt.Sprintf("%v:%v", ip, port)
	// return values from the proble... one should be parsed... maybe? for seconds.  will do that later.
	out, out2 := k.ExecWithOptions([]string{"curl", pod.Status.PodIP}, ns1, pod1, "busybox")
	if err == nil {
		fmt.Println("success", out, out2)
		return 5
	}
	return -1
}

// ExecWithOptions executes a command in the specified container,
// returning stdout, stderr and error. `options` allowed for
// additional parameters to be passed.
func (k *Kubernetes) ExecWithOptions(Command []string, ns string, pod string, cname string) (string, string) {
	fmt.Println("Exec...")

	req := k.Client().RESTClient().Post().
		Resource("pods").
		Name(pod).
		Namespace(ns).
		SubResource("exec").
		Param("container", cname)

	req.VersionedParams(&v1.PodExecOptions{
		Container: cname,
		Command:   Command,
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}, scheme.ParameterCodec)

	var stdout, stderr bytes.Buffer
	fmt.Println(req.URL())
	err := execute("POST", req.URL(), kubeconfig(), nil, &stdout, &stderr, true)
	if err != nil {
		fmt.Println("error: ", err, " ,,,, ", stdout.String(), stderr.String())
		panic("failing...")
	}
	return strings.TrimSpace("out=" + stdout.String()), "err=" + strings.TrimSpace(stderr.String())
}

func execute(method string, url *url.URL, config *restclient.Config, stdin io.Reader, stdout, stderr io.Writer, tty bool) error {
	exec, err := remotecommand.NewSPDYExecutor(config, method, url)
	if err != nil {
		return err
	}
	return exec.Stream(remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
		Tty:    tty,
	})
}

func kubeconfig() *rest.Config {
	// TODO borrowed from e2e upstream.
	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	return config
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
			Name:   n,
			Labels: labels,
		},
	}
	nsr, err := k.Client().CoreV1().Namespaces().Create(ns)
	fmt.Println(err)
	return nsr, err
}

func (k *Kubernetes) CreateDeployment(ns, deploymentName string, replicas int32, labels map[string]string, image string) (*appsv1.Deployment, error) {
	zero := int64(0)
	fmt.Println("ns", ns)
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
					Labels:    labels,
					Namespace: ns,
				},
				Spec: v1.PodSpec{
					TerminationGracePeriodSeconds: &zero,
					Containers: []v1.Container{
						{
							Name:            image,
							Image:           image,
							SecurityContext: &v1.SecurityContext{},
							Command:         []string{"sleep", "600"},
						},
					},
				},
			},
		},
	}

	d, err := k.Client().AppsV1().Deployments(ns).Create(d)
	fmt.Println(err)
	return d, err
}

func (k *Kubernetes) CreateNetworkPolicy(netpol *v1net.NetworkPolicy) (*v1net.NetworkPolicy, error) {
	np, err := k.Client().NetworkingV1().NetworkPolicies("").Create(netpol)
	return np, err
}
