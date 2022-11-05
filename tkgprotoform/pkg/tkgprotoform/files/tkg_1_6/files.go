package tkg_1_6

import (
	_ "embed"
)

// package up the raw input files using go embed.
// these get hydrated at runtime w/ parameters
var (
	//go:embed cluster.yaml
	TKG_16_ClusterConfig string

	//go:embed 2_management_cluster.sh
	TKG_16_ManagementClusterInstall string

	// //go:embed image_builder.sh
	TKG_16_ImageBuilderScript string
)

// files has all the files we are going to write to disk, so that
// they can be hacked up by the user after running tkgprotoform init the first time.
func Files() map[string]string {
	return map[string]string{
		"cluster.yaml":            TKG_16_ClusterConfig,
		"2_management_cluster.sh": TKG_16_ManagementClusterInstall,
		"image_builder.sh":        TKG_16_ImageBuilderScript,
	}
}
