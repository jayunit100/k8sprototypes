This assumes you have a linux kube proxy running already and only windows (kube proxy, cni) are missing.
This just puts everything in `kube-system` . Some stuff is hardcoded to kube-system.

# First

Install calico-main.yaml .  It has hacks that run the "linux" calico stuff in a way that they wont die when windows calico agents dont implement BGP.  i.e. it removes bird liveness probes.

# Next

Now install the kube-proxy yaml.  That will run as host process container.  Modify the `image:`  to docker.io if you need to.

# Finally

Now you can run calico on *windows*.  Install the calico-windows yaml file.  Easy.  

Now what?  Your calico pods might not be happy b/c HNS network not created by the calico host process pod. 

# Fix the HNS network

This is hacky but - just run the ps1 script - SSH into your nodes one by one and run it.  If this is TKG it'll work ootb, and fail w/ an 
error realted to CNI dirs.  Thats fine.  All you care is that it got far enough to create the HNS network.  YOu'll know it worked bc after you run it
you'll see that your kube proxy comes up , (because HNS is now available).

And right after **kube proxy** comes online, youll see the calico pods come up!
