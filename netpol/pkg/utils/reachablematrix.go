package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

// TODO A Graph data model might be better for this class.
type ReachableMatrix struct {
	DefaultExpect bool
	Pods          []string
	Namespaces    []string
	Expected      map[string]map[string]bool
	Observed      map[string]map[string]bool
}

func addToMap(m map[string]map[string]bool, ns string, pod string, ns2 string, pod2 string, connectionTime bool) {
	if m == nil {
		fmt.Println("init")
		m = map[string]map[string]bool{}
	}
	if m[ns+"_"+pod] == nil {
		m[ns+"_"+pod] = map[string]bool{}
	}
	m[ns+"_"+pod][ns2+"_"+pod2] = connectionTime
}

func (r *ReachableMatrix) init() {
	if r.Expected != nil {
		return
	}

	// create datastructures if not created yet.
	r.Expected = map[string]map[string]bool{}
	r.Observed = map[string]map[string]bool{}

	if r.Pods == nil || r.Namespaces == nil {
		panic("pods/ns must be existent in the struct.")
	}
	conn := r.DefaultExpect

	// create the matrix
	for _, n := range r.Namespaces {
		for _, p := range r.Pods {
			for _, nn := range r.Namespaces {
				for _, pp := range r.Pods {
					addToMap(r.Expected, n, p, nn, pp, conn)
				}
			}
		}
	}
}

func (r *ReachableMatrix) HasNS(n, pod string) bool {
	for _, v := range r.Namespaces {
		if v == n {
			return true
		}
	}
	return false
}

func (r *ReachableMatrix) Expect(ns string, pod string, ns2 string, pod2 string, conn bool) {
	r.init()
	if r.HasNS(ns, pod) && r.HasNS(ns2, pod2) {
		addToMap(r.Expected, ns, pod, ns2, pod2, conn)
	} else {
		panic(fmt.Sprintf("ns/pod not found %v %v ", ns, pod))
	}
}

// ExpectAllIngress defines that any traffic going into the pod in 'ns' will be allowed/denied (true/false)
func (r *ReachableMatrix) ExpectAllIngress(ns, pod string, connected bool) {
	r.init()

	for _, nsFrom := range r.Namespaces {
		for _, podFrom := range r.Pods {
			r.Expect(nsFrom, podFrom, ns, pod, connected)
			if !connected {
				log.Infof("Blacklisting %v %v %v %v", nsFrom, podFrom, ns, pod)
			}
		}
	}
}

func (r *ReachableMatrix) Observe(ns string, pod string, ns2 string, pod2 string, conn bool) {
	r.init()
	if r.HasNS(ns, pod) && r.HasNS(ns2, pod2) {
		addToMap(r.Observed, ns, pod, ns2, pod2, conn)
	} else {
		panic(fmt.Sprintf("ns/pod not found %v %v ", ns, pod))
	}
}

func (r *ReachableMatrix) GetExpectedObserved(ns string, pod string, ns2 string, pod2 string) (bool, bool) {
	r.init()
	return r.Expected[ns+"_"+pod][ns2+"_"+pod2], r.Observed[ns+"_"+pod][ns2+"_"+pod2]
}

func (r *ReachableMatrix) LengthExpectedObserved() (int, int) {
	r.init()
	return len(r.Expected), len(r.Observed)
}

func Key(ns, pod string) string {
	return ns + "_" + pod
}

func (r *ReachableMatrix) Summary() (string, bool) {
	falseObs := 0
	trueObs := 0
	for _, n1 := range r.Namespaces {
		for _, p1 := range r.Pods {
			for _, n2 := range r.Namespaces {
				for _, p2 := range r.Pods {
					from := Key(n1, p1)
					to := Key(n2, p2)
					if _, ok := r.Observed[from][to]; !ok {
						log.Infof("WARNING ----> Observation vals not done yet... from:%v, to:%v obsFromPods:%v, matrix:%v", from, to, len(r.Observed[from]), r.Observed[from])
					}
					if r.Expected[from][to] == r.Observed[from][to] {
						trueObs++
					} else {
						//fmt.Print(from, "->", to, " not matching expect=")
						//fmt.Print(r.Expected[from][to], ", observed=")
						//fmt.Println(r.Observed[from][to])
						falseObs++
					}
				}
			}
		}
	}
	for fromPod, dict := range r.Observed {
		fmt.Println(fromPod)
		line := []string{}
		for toPod, v := range dict {
			val := "F"
			if v {
				val = "T"
			}
			line = append(line, fmt.Sprintf("%s:%s", toPod, val))
		}
		fmt.Println("-->", line)
	}
	passed := falseObs == 0
	return fmt.Sprintf("correct:%v, incorrect:%v, result=%t", trueObs, falseObs, passed), passed
}
