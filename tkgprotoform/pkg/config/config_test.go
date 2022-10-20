package config

import "testing"

func TestConf(t *testing.T) {

	c := GetConf()

	if len(c.TanzuInputs) == 0 {
		t.Logf("tanzu inputs ... %v ... is wrong", c.TanzuInputs)
		t.Fail()
	}
}
