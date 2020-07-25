# CNI Debugging 

How can I debug the datapath of containers in my Kubernetes clusters using basic linux commands? 


The first step is to run `ip a`.  This will give you a birds eye view of what network interfaces your host is aware of.

## First, look at interfaces with the IP command.

In our calico cluster, running `ip a` shows us that we have a `tunl0` interface.  This is created by the calico/node container,
via the brd service which is responsible for routing traffic through the IPIP tunnel in the cluster.

```
# ip a
2: tunl0@NONE: <NOARP,UP,LOWER_UP> mtu 1440 qdisc noqueue state UNKNOWN group default qlen 1000
```

In antrea... 

```
capv@antrea-vc-md-0-869559ff88-6mnd4 [ ~ ]$ ip a
3: ovs-system: <BROADCAST,MULTICAST> mtu 1500 qdisc noop state DOWN group default qlen 1000
    link/ether 7e:de:21:4b:88:46 brd ff:ff:ff:ff:ff:ff
5: antrea-gw0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1450 qdisc noqueue state UNKNOWN group default qlen 1000
    link/ether 82:aa:a9:6f:02:33 brd ff:ff:ff:ff:ff:ff
    inet 100.96.29.1/24 brd 100.96.29.255 scope global antrea-gw0
       valid_lft forever preferred_lft forever
    inet6 fe80::80aa:a9ff:fe6f:233/64 scope link
```

In both clusters, run `kubectl scale deployment coredns --replicas=10 -n kube-system`.  Then re run these commands.
You'll notices that you do not see any new ip.

### How many packets are flowing through the interfaces? 

The `ip` command actually has a `-s` option, which will show you if traffic is flowing...

```
10: cali3317e4b4ab5@if5: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1440 qdisc noqueue state UP group default
    link/ether ee:ee:ee:ee:ee:ee brd ff:ff:ff:ff:ff:ff link-netns cni-abb79f5f-b6b0-f548-3222-34b5eec7c94f
    RX: bytes  packets  errors  dropped overrun mcast
    150575     1865     0       2       0       0
    TX: bytes  packets  errors  dropped carrier collsns
    839360     1919     0       0       0       0
```

This command works the same in antrea, however, we run it against the gatweay bridge to see total amount of traffic.
```
5: antrea-gw0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1450 qdisc noqueue state UNKNOWN group default qlen 1000
    link/ether 82:aa:a9:6f:02:33 brd ff:ff:ff:ff:ff:ff
    inet 100.96.29.1/24 brd 100.96.29.255 scope global antrea-gw0
       valid_lft forever preferred_lft forever
    inet6 fe80::80aa:a9ff:fe6f:233/64 scope link
       valid_lft forever preferred_lft forever
    RX: bytes  packets  errors  dropped overrun mcast
    89662090   1089577  0       0       0       0
    TX: bytes  packets  errors  dropped carrier collsns
    108901694  1208573  0       0       0       0
```

# Route

The next level of debugging involves seeing how these devices are wired to IP addresses.  You might need to run `apt-get update; apt-get install net-tools` in your kind cluster
if you dont have a full blown linux vm setup. 

Running `route -n` in our calico cluster shows the following routing table in the kernel:

```
	# route -n
	Kernel IP routing table
	Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
	0.0.0.0         172.18.0.1      0.0.0.0         UG    0      0        0 eth0
	172.18.0.0      0.0.0.0         255.255.0.0     U     0      0        0 eth0
	192.168.9.128   172.18.0.3      255.255.255.192 UG    0      0        0 tunl0
	192.168.71.0    172.18.0.5      255.255.255.192 UG    0      0        0 tunl0
	192.168.88.0    172.18.0.4      255.255.255.192 UG    0      0        0 tunl0
	192.168.143.64  172.18.0.2      255.255.255.192 UG    0      0        0 tunl0
	192.168.173.64  0.0.0.0         255.255.255.192 U     0      0        0 *
	192.168.173.65  0.0.0.0         255.255.255.255 UH    0      0        0 calicd2f389598e
	192.168.173.66  0.0.0.0         255.255.255.255 UH    0      0        0 calibaa5769d671
```

What about our bridge based CNI?  

Interestingly, we don't see a new destination IP for every device.  Instead, we see that there is a `.1` antrea-gateway.

```
	root [ /home/capv ]# route -n
	Kernel IP routing table
	Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
	0.0.0.0         192.168.5.1     0.0.0.0         UG    1024   0        0 eth0
	100.96.0.0      100.96.0.1      255.255.255.0   UG    0      0        0 antrea-gw0
	100.96.21.0     100.96.21.1     255.255.255.0   UG    0      0        0 antrea-gw0
	100.96.26.0     100.96.26.1     255.255.255.0   UG    0      0        0 antrea-gw0
	100.96.28.0     100.96.28.1     255.255.255.0   UG    0      0        0 antrea-gw0
```

Differences:
- The bitmask is /24, which means ALL IP addresses destined for our 1st node go to 100.96.0.0
- All IP addresses destined for our next node go to 100.96.21.0,
- All IP addresses destined for the next node go do 100.96.26.0

The differences here is
- Antrea has one routing table entry PER NODE
- Calico has one routing table entry PER POD

### Destinations and gateways for calico

The calico cluster (which has pods on the 192 subnet) interestingly has *Gateway* IP addresses on a *Different* subnet.

In a multinode calico cluster, the `route -n` command showed us that the *gateway* IP is actually the same as that
of the docker node that a pod lives on .  This tells us that packets can be routed directly to nodes, and that the
nodes have the ability to route traffics directly from the host.  If we look closely at the calico processes, we'll see 
indeed that the route table has a specific route for each and every calico container.
```
	192.168.9.128   172.18.0.3      255.255.255.192 UG    0      0        0 tunl0
	192.168.71.0    172.18.0.5      255.255.255.192 UG    0      0        0 tunl0
	192.168.88.0    172.18.0.4      255.255.255.192 UG    0      0        0 tunl0
```

This is because calico will route a pod to a node IP, as the gateway, whereas, antrea will route a pod to the
antrea-gw device which is a process that is accessible over the pod subnet as we'll see when we look at the antrea table's
Destination and Gatway fields.

### Destinations and Gateways in the antrea subnets

The *destination* and *gateway* fields are obviously the most relevant here.  

- The *Destination* field tells us where a packet is heading, for example, a packet heading
to 100.96.21.13 would end up in the 3rd row above.
- The *Gateway* field tells us where a packet going to go *before* it reaches the destination, i.e.
the *next hop* in its path to its ultimate destination.  In this case the 100.96.21.13 packet
would be headed to *100.96.21.1*.

In calico, the gateway for all packets is 0.0.0.0, and the interfaces are explicit for every pod.
In antrea, the gateway ip addresses end in .1.  Lets see what IP that corresponds to...

On this node, we can see that 192.96.29.0 maps to the local address, but 100.96.26.1 maps to a *different* address.
```
        100.96.26.0     100.96.26.1
	100.96.29.0     0.0.0.0
```
Interestingly, the *gateway* ip address 100.96.26.1 which is going to be the destination for all .26 traffic isnt the
address of a pod but is on the pod subnet (100).  Where does it come from? 

And indeed,  we ran `ip a` earlier, the ip address 100.96.29.1 is associated with the `antrea-gw0` device.
Thus, that device is the *gateway* for ALL traffic.  Every antrea container will create a gateway device and bind it to this
.1 ip on the subnet which it owns.

We can then ask open vswitch to tell us more about this interface, via the `ovs-vsctl` tool:

```
root@antrea-vc-control-plane-6njtk:/# ovs-vsctl list interface|grep -A 5 antrea
name                : antrea-gw0
ofport              : 2
ofport_request      : 2
options             : {}
other_config        : {}
statistics          : {collisions=0, rx_bytes=1773391201, rx_crc_err=0, rx_dropped=0, rx_errors=0, rx_frame_err=0, rx_missed_errors=0, rx_over_err=0, rx_packets=16392260, tx_bytes=6090558410, tx_dropped=0, tx_errors=0, tx_packets=17952545}
```

Note that to do this, we executed the command directly inside of the `antrea-agent` command.

## TCPDump

Now that we have traced the relationship between our host to the underlying linux networking tools which route traffic,
we may want to look at things from the container's perspective.

The most common tool for doing this is `tcpdump`.  Lets grab one of our coredns containers and look at its traffic.

In calico we can directly sniff the packets on the `cali` devices, like so:
```
192.168.173.66  0.0.0.0         255.255.255.255 UH    0      0        0 calibaa5769d671
# tcpdump -i calicd2f389598e
tcpdump: verbose output suppressed, use -v or -vv for full protocol decode
listening on calicd2f389598e, link-type EN10MB (Ethernet), capture size 262144 bytes
20:13:07.733139 IP 10.96.0.1.443 > 192.168.173.65.60684: Flags [P.], seq 1615967839:1615968486, ack 1173977013, win 264, options [nop,nop,TS val 296478
```

The 10.96.0.1 ip address is the internal kubernetes address.  Its acknlowedging receipt of a request from the coredns server to get a dns record.

If we look at a typical node in our cluster where we're running the coredns pod, our antrea pods will be named like so:
```
30: coredns--e5cc00@if3: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1450 qdisc noqueue master ovs-system state UP group default
    link/ether e6:8a:27:05:d7:30 brd ff:ff:ff:ff:ff:ff link-netns cni-2c6b1bc0-cf36-132c-dfcb-88dd158f51ca
    inet6 fe80::e48a:27ff:fe05:d730/64 scope link
       valid_lft forever preferred_lft forever
```

This means we can directly sniff the packets going to this node by attaching to this veth device with tcpdump:

```
tcpdump -i coredns--29244a -n
```

Be sure to remove the `@f3` part when sending this .  When you run this command, you should see traffic from different pods which
are attempting to resolve kubernetes DNS records.

We often use the -n option so that ip addresses don't get hidden from us when using tcpdump.

If you specifically want to see if one pod is talking to another, you can go to the node on the pod you are receiving traffic on, and scrape all tcp traffic which inclues one of the Pod's ip addresses.  Say that pod that is sending traffic is 100.96.21.21:

```
root [ /home/capv ]# tcpdump host 100.96.21.21 -i coredns--29244a
tcpdump: verbose output suppressed, use -v or -vv for full protocol decode
listening on coredns--29244a, link-type EN10MB (Ethernet), capture size 262144 bytes
21:59:36.818933 IP 100.96.21.21.45978 > 100.96.26.19.9153: Flags [S], seq 375193568, win 64860, options [mss 1410,sackOK,TS val 259983321 ecr 0,nop,wscale 7], length 0
21:59:36.819008 IP 100.96.26.19.9153 > 100.96.21.21.45978: Flags [S.], seq 3927639393, ack 375193569, win 64308, options [mss 1410,sackOK,TS val 2440057191 ecr 259983321,nop,wscale 7], length 0
21:59:36.819928 IP 100.96.21.21.45978 > 100.96.26.19.9153: Flags [.], ack 1, win 507, options [nop,nop,TS val 259983323 ecr 2440057191], length 0
2
```

This will give you a raw dump of anything with, for example, a 192 address, and a 9153 port.

## IPTables

Once you've looked at wether individual pods are or are not receiving traffic from other pods in your cluster 

























































