package utils

import "fmt"

type ReachableMatrix struct {
	DefaultExpect bool
	Pods          []string
	Namespaces    []string
	Expected      map[string]map[string]bool
	Observed      map[string]map[string]bool
}

func addToMap(m map[string]map[string]bool, ns string, pod string, ns2 string, pod2 string, connectionTime bool) {
	if m == nil {
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
	if r.Pods == nil || r.Namespaces == nil {
		panic("pods/ns must be existent in the struct.")
	}
	conn := false
	if r.DefaultExpect == true {
		conn = true // expect connect within 5 seconds
	}
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

func (r *ReachableMatrix) Expect(ns string, pod string, ns2 string, pod2 string, conn bool) {
	r.init()

	addToMap(r.Expected, ns, pod, ns2, pod2, conn)
}

func (r *ReachableMatrix) Observe(ns string, pod string, ns2 string, pod2 string, conn bool) {
	r.init()
	addToMap(r.Observed, ns, pod, ns2, pod2, conn)
}

func (r *ReachableMatrix) GetExpectedObserverd(ns string, pod string, ns2 string, pod2 string) (bool, bool) {
	r.init()
	return r.Expected[ns+"_"+pod][ns2+"_"+pod2], r.Observed[ns+"_"+pod][ns2+"_"+pod2]
}

func (r *ReachableMatrix) LengthExpectedObserved() (int, int) {
	r.init()
	return len(r.Expected), len(r.Observed)
}

func (r *ReachableMatrix) Summary() (string, bool) {
	falseObs := 0
	trueObs := 0
	for _, n := range r.Namespaces {
		for _, p := range r.Pods {
			if r.Expected[n][p] == r.Observed[n][p] {
				trueObs++
			} else {
				falseObs++
			}
		}
	}
	passed := falseObs == 0
	return fmt.Sprintf("correct:%v, incorrect:%v, result=", trueObs, falseObs, passed), passed
}
