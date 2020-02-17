package utils

import "fmt"

// TODO A Graph data model might be better for this class.
type ReachableMatrix struct {
	DefaultExpect bool
	Pods          []string
	Namespaces    []string
	Expected      map[string]map[string]bool
	Observed      map[string]map[string]bool
}

func addToMap(m map[string]map[string]bool, ns string, pod string, ns2 string, pod2 string, connectionTime bool) {
	fmt.Println(" map length", len(m))

	if m == nil {
		fmt.Println("init")
		m = map[string]map[string]bool{}
	}
	if m[ns+"_"+pod] == nil {
		m[ns+"_"+pod] = map[string]bool{}
	}
	m[ns+"_"+pod][ns2+"_"+pod2] = connectionTime
	fmt.Println("new map length", len(m), m)
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

func (r *ReachableMatrix) Observe(ns string, pod string, ns2 string, pod2 string, conn bool) {
	r.init()
	if r.HasNS(ns, pod) && r.HasNS(ns2, pod2) {
		addToMap(r.Observed, ns, pod, ns2, pod2, conn)
	} else {
		panic(fmt.Sprintf("ns/pod not found %v %v ", ns, pod))
	}
	fmt.Println("observing ", ns, pod, ns2, pod2, conn, " new len ", len(r.Observed))
}

func (r *ReachableMatrix) GetExpectedObserverd(ns string, pod string, ns2 string, pod2 string) (bool, bool) {
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
					fmt.Println("checking that all expected and observed values are written.", n1, p1, n2, p2)
					fmt.Println(len(r.Observed), len(r.Expected))
					if _, ok := r.Expected[Key(n1, p1)][Key(n2, p2)]; !ok {
						panic("missing Expected val")
					}
					if _, ok := r.Observed[Key(n1, p1)][Key(n2, p2)]; !ok {
						panic("missing Observation val")
					}

					if r.Expected[Key(n1, p1)][Key(n2, p2)] == r.Observed[Key(n1, p1)][Key(n2, p2)] {
						trueObs++
					} else {
						fmt.Print(n1, p1, "->", n2, p2, " not matching expect=")
						fmt.Print(r.Expected[Key(n1, p1)][Key(n2, p2)], ", observed=")
						fmt.Println(r.Observed[Key(n1, p1)][Key(n2, p2)])
						falseObs++
					}
				}
			}
		}
	}
	passed := falseObs == 0
	return fmt.Sprintf("correct:%v, incorrect:%v, result=", trueObs, falseObs, passed), passed
}
