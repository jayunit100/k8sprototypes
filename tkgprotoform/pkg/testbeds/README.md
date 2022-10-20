#  Testbeds

There are different "test beds" which the various tanzu operations must run on.

One such testbed is the "Empty" testbed, i.e. your local laptop, installing Tanzu.

However even your local laptop requires some bootstrapping:

- You need a payload directory where the tanzu tarballs and OVAs will live
- You need to run tanzu init and download the appropriate version / confirm it runs
- You need to verify docker is installed and so on.

Each "Testbed" implementation should "set up" whatever the basic compute infra is 
in order for a user to begin running regular commands like:

```
tanzu mc create -f my-cluster.yaml
```

That means:

- Setting up an NSX topology
- Setting up any EC2 or AZure pre requisites/verifying cloud keys are valid
- Creating a VSphere and or NSX data center w/ an IP address
- ... 
