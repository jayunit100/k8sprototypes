package tkg_1_6

import (
	_ "embed"
)

// package up the raw input files using go embed.
// these get hydrated at runtime w/ parameters
var (
	//go:embed cluster.yaml
	TKG_16_ClusterConfig string

	// Setup of TKG installation artifacts...

	//go:embed 0_install_payload.sh
	TKG_16_Install_Payload string

	//go:embed 0_verify_payload.sh
	TKG_16_Verify_Payload string

	//go:embed 1_upload_ova.sh
	TKG_16_UploadOVA string

	//go:embed 2_management_cluster.sh
	TKG_16_ManagementClusterInstall string

	//go:embed 2_swap_tkr.sh
	TKG_16_Swap_TKR string

	//go:embed 3_workload_cluster.sh
	TKG_16_WorkloadClusterInstall string

	// //go:embed image_builder.sh
	TKG_16_ImageBuilderScript string

	//go:embed image-builder.json
	TKG_16_ImageBuilderJSON string

	//go:embed image-builder-credentials.json
	TKG_16_ImageBuilderCredsJSON string

	//go:embed image-builder-customizations.json
	TKG_16_ImageBuilderCustJSON string

	//go:embed image-builder-tkg.json
	TKG_16_ImageBuilderTKGJSON string
)

// files has all the files we are going to write to disk, so that
// they can be hacked up by the user after running tkgprotoform init the first time.
func Files() map[string]string {
	return map[string]string{

		// Build Cluster images, may be used for mgmt or wl ...
		"image_builder.sh":                  TKG_16_ImageBuilderScript,
		"image-builder.json":                TKG_16_ImageBuilderJSON,
		"image-builder-credentials.json":    TKG_16_ImageBuilderCredsJSON,
		"image-builder-customizations.json": TKG_16_ImageBuilderCustJSON,
		"image-builder-tkg.json":            TKG_16_ImageBuilderTKGJSON,

		"cluster.yaml": TKG_16_ClusterConfig,

		// Install Tanzu CLI
		"0_install_payload.sh": TKG_16_Install_Payload,
		"0_verify_payload.sh":  TKG_16_Verify_Payload,

		// Initialize TKG, generic steps
		"1_upload_ova.sh":         TKG_16_UploadOVA,
		"2_management_cluster.sh": TKG_16_ManagementClusterInstall,

		// make a new TKRs (optional) and make WL cluster
		"2_swap_tkr.sh": TKG_16_Swap_TKR,

		// Finally, create WL clusters...
		"3_workload_cluster.sh": TKG_16_WorkloadClusterInstall,
	}
}
