package utils

type ReachableMatrix struct {
	Expected map[string]map[string]int
	Observed map[string]map[string]int
}

func addToMap(m map[string]map[string]int, ns string, pod string, ns2 string, pod2 string, connectionTime int) {
	if m == nil {
		m = map[string]map[string]int{}
	}
	if m[ns+"_"+pod] == nil {
		m[ns+"_"+pod] = map[string]int{}
	}
	m[ns+"_"+pod][ns2+"_"+pod2] = connectionTime
}

func (r *ReachableMatrix) Expect(ns string, pod string, ns2 string, pod2 string, connectionTime int) {
	addToMap(r.Expected, ns, pod, ns2, pod2, connectionTime)
}

func (r *ReachableMatrix) Observe(ns string, pod string, ns2 string, pod2 string, connectionTime int) {
	addToMap(r.Observed, ns, pod, ns2, pod2, connectionTime)
}

func (r *ReachableMatrix) GetExpectedObserverd(ns string, pod string, ns2 string, pod2 string) (int, int) {
	return r.Expected[ns+"_"+pod][ns2+"_"+pod2], r.Observed[ns+"_"+pod][ns2+"_"+pod2]
}

func (r *ReachableMatrix) LengthExpectedObserved() (int, int) {
	return len(r.Expected), len(r.Observed)
}
