# Prototypical installation of Windows on TKG

This is not supported by VMWare at this time, but rather a place to iterate on upstream artifacts for windows support on Cluster API with the Cluster APIVSphere provider.  Some of the artifacts referenced may require having access to a VMWare provided K8s installation, but they can all be built from equivalent upstream components.

# Step 0: Instal a Management cluster

You need a CAPV managmement cluster to start

# Step 1: Create a Burrito Service

The burrito service will host K8s binaries on an endpoint for you, it will run inside 
of your K8s cluster.

# Step 2: Run Image builder to create an OVA inside a K8s management cluster

The image builder tool will pull down artifacts from your running burrito service.
The image builder tool will create an OVA for you, and upload it to Vcenter.

# Step 3:

USe the YTT Templates in this repository to make a cluster using Tanzu CLI.
