package config

import "testing"

func TestConf(t *testing.T) {

	c := GetConf()

	if len(c.TanzuInputs) == 0 {
		t.Logf("tanzu inputs ... %v ... is wrong", c.TanzuInputs)
		t.Fail()
	} else {
		// We expect this to have VSPHERE_CONTROL_PLANE_ENDPOINT, bc that is
		// a valid input to the config.yaml
		t.Log("Acquired inputs", c.TanzuInputs)
	}
}
