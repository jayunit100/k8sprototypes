# What does a CAPV recovery scenario look like ? 

In CAPV HAProxy can be used to indicate the health of yourclusters.

Building off of what we know from how etcd operates in duress, lets consider some experiments

and see how the ClusterAPI responds... 

In a running 3 node CAPV cluster, deleting one etcd node wont cause writes to block.
But, deleting 2 nodes will !

In a running CAPV cluster, we can delete a few nodes:

![Image description](etcd_starvation_after_killing_node.png)

After doing this we can power back on machines, until we restore 2/3 quorum:

![Image description](etcd_backonline_after_powering_1_node.png)

Finally, we can manually power back on the 3rd node...

![Image description](fully_recovered_etcd.png)
