package config

import (
	"testing"
	"tkgprotoform.com/protoform/pkg/util"
)

func TestConf(t *testing.T) {

	c := GetConf()

	if util.FileExists("config.yaml") {
		t.Log("Found config.yaml")
	} else {
		t.Log("Missing config.yaml")
		t.Fail()
	}
	// This directory has a "config.yaml" in it.  Thus,
	// we expect
	if len(c.TanzuInputs) == 0 {
		t.Logf("tanzu inputs ... %v ... is wrong", c.TanzuInputs)
		t.Fail()
	} else {
		// We expect this to have VSPHERE_CONTROL_PLANE_ENDPOINT, bc that is
		// a valid input to the config.yaml
		t.Log("Acquired inputs", c.TanzuInputs)
	}
}
