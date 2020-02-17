package utils

import (

	//	"context"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	v1net "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

type Kubernetes struct {
}

func (k *Kubernetes) GetPods(ns string, key, val string) []v1.Pod {
	v1PodList, _ := Client().CoreV1().Pods(ns).List(metav1.ListOptions{
		LabelSelector: "pod",
	})
	pods := []v1.Pod{}
	for _, pod := range v1PodList.Items {
		if pod.Labels[key] == val {
			pods = append(pods, pod)
		}
	}
	return v1PodList.Items
}

func (k *Kubernetes) Probe(ns1 string, pod1 string, ns2 string, pod2 string, port int) bool {
	ip := "1.1.1.1"
	pod, err := Client().CoreV1().Pods(ns2).Get(pod2, metav1.GetOptions{})
	if err == nil {
		ip = pod.Status.PodIP
	} else {
		panic(err)
	}
	//fmt.Println("Pod ip", pod.Status.PodIP)
	// delete index.html before curling.

	//fmt.Println("exec starts now ...")
	_, _, err = ExecuteRemoteCommand(pod, []string{"rm", "-f", "index.html"})
	out, out2, err := ExecuteRemoteCommand(pod, []string{"wget", "http://" + ip + ":" + fmt.Sprintf("%v", port)})
	if err == nil {
		//	fmt.Println("success", "out="+out, "err="+out2)
		return true
	} else {
		fmt.Println("\n\n", out, out2, "\n\n")
		panic(err)
	}
	return false
}

// ExecuteRemoteCommand executes a remote shell command on the given pod
// returns the output from stdout and stderr
func ExecuteRemoteCommand(pod *v1.Pod, command []string) (string, string, error) {
	kubeCfg := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)
	restCfg, err := kubeCfg.ClientConfig()
	if err != nil {
		return "", "", err
	}
	//coreClient, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		return "", "", err
	}

	buf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	request := Client().CoreV1().RESTClient().Post().Namespace(pod.Namespace).Resource("pods").
		Name(pod.Name).SubResource("exec").VersionedParams(&v1.PodExecOptions{
		Command: command,
		Stdin:   false,
		Stdout:  true,
		Stderr:  true,
		TTY:     true},
		scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(restCfg, "POST", request.URL())
	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: buf,
		Stderr: errBuf, ///home/jayunit100/go/src/github.com/jayunit100/k8sprototypes/netpol/pkg/utils/k8s_util.gohome/jayunit100/go/src/github.com/jayunit100/k8sprototypes/netpol/pkg/utils/k8s_util.go/home/jayunit100/go/src/github.com/jayunit100/k8sprototypes/netpol/pkg/utils/k8s_util.goome/jayunit100/go/src/github.com/jayunit100/k8sprototypes/netpol/pkg/utils/k8s_util.go
	})
	if err != nil {
		return buf.String(), errBuf.String(), errors.Wrapf(err, "Failed executing command %s on %v/%v------/%v/%v", command, pod.Namespace, pod.Name, buf.String(), errBuf.String())
	}
	return buf.String(), errBuf.String(), nil
}

func Client() *kubernetes.Clientset {
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
	nsr, err := Client().CoreV1().Namespaces().Create(ns)
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
							Name:            "prober",
							Image:           image,
							SecurityContext: &v1.SecurityContext{},
							// Command:         []string{"sleep", "60000"},
							Ports: []v1.ContainerPort{
								v1.ContainerPort{
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	d, err := Client().AppsV1().Deployments(ns).Create(d)
	fmt.Println(err)
	return d, err
}

func (k *Kubernetes) CreateNetworkPolicy(netpol *v1net.NetworkPolicy) (*v1net.NetworkPolicy, error) {
	np, err := Client().NetworkingV1().NetworkPolicies("").Create(netpol)
	return np, err
}
