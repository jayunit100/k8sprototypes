package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type Test struct {
	focus string `json:"focus"`
	skip string `json:"skip"`
}

func ReadTests() map[string]map[string]string {

	bytes := []byte(`{"Core":{ "focus":"windows", "skip": "gmsa"},"Networking":{ "focus":"NetworkPolicy", "skip": "6|udp" }}`)
	var dat map[string]map[string]string

	if err := json.Unmarshal(bytes, &dat); err != nil {
		panic(err)
	}
	fmt.Println(dat)
	return dat
}


func runTest(t map[string]string) (string, error){
	args := []string{
		"--ginkgo.v=true",
		"--ginkgo.debug=true",
		"--kubeconfig=/home/kubo/.kube/config",
		"--ginkgo.focus=%v",
		"--node-os-distro=windows",
		"--ginkgo.skip=%v",
		"--ginkgo.noColor=true",
		"--non-blocking-taints=\"os,node-role.kubernetes.io/master,node.kubernetes.io/not-ready\"",
	}
	argsUsed := fmt.Sprintf(strings.Join(args, " "), t["focus"], t["skip"])

	split := strings.Split(argsUsed, " ")

	// date
	fmt.Println(argsUsed)
	runme := exec.Command("./e2e.test", split...)
	out, err := runme.CombinedOutput()
	return string(out), err

}

func main() {
	tests := ReadTests()

	for n,t := range tests {
		s := fmt.Sprintf("Operational Readiness: %v", n)
		fmt.Println(s)
		o, e := runTest(t)
		fmt.Println(o)
		fmt.Println(e)
	}
}