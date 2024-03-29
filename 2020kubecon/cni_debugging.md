# CNI Debugging 

https://www.youtube.com/watch?v=RQNy1PHd5_A

How can I debug the datapath of containers in my Kubernetes clusters using basic linux commands? 

- Rough outline of pod -> service -> node -> policy flow, generalized for any CNI.
```
+-----------------------------+    +----------------------------+
|                             |    |                            |
|               100..61       |    |                            |
| +--------+      +-------+   |    |  +---------+     +------+  |
| |   pod  +------+service|   |    |  |  policy +---->+ pod  |  |
| +--------+      +-+-----+   |    |  +----^----+     +------+  |
|  100.96           |         |    |       |              100.96|
|               +---v------+  |    |       |                    |
|               |  iptables+----+  | +-----+----------------+   |
|               +----------+  | +---->node firewall: tcpdump|   |
|                             |    | +----------------------+   |
|                             |    |                            |
+-----------------------------+    +----------------------------+
                 10...                          10...
```
The first step is to run `ip a`.  This will give you a birds eye view of what network interfaces your host is aware of.

In terms of the IPTables routing, we have the following:
```
KUBE_MARK_MASQ -> KUBE-SVC -----> KUBE_MARK_DROP  
			  |-----> KUBE_SEP -> KUBE_MARQ_MASK -> NODE -> route device) 
```

## Quick over view of eth pairs

Any k8s cluster node will have a list of containers mapped to mac addresses ...  
- Addresses of pods that the node can reach:
```
root [ /home/capv ]# arp -na | sort
? (100.96.26.15) at 86:55:7a:e3:73:71 [ether] on antrea-gw0
? (100.96.26.16) at 4a:ee:27:03:1d:c6 [ether] on antrea-gw0
? (100.96.26.17) at <incomplete> on antrea-gw0
? (100.96.26.18) at ba:fe:0f:3c:29:d9 [ether] on antrea-gw0
? (100.96.26.19) at e2:99:63:53:a9:68 [ether] on antrea-gw0
? (100.96.26.20) at ba:46:5e:de:d8:bc [ether] on antrea-gw0
? (100.96.26.21) at ce:00:32:c0:ce:ec [ether] on antrea-gw0
? (100.96.26.22) at e2:10:0b:60:ab:bb [ether] on antrea-gw0
? (100.96.26.2) at 1a:37:67:98:d8:75 [ether] on antrea-gw0
```
- Addresses which are local to the node:
```
? (192.168.5.160) at 00:50:56:b0:ee:ff [ether] on eth0
? (192.168.5.1) at 02:50:56:56:44:52 [ether] on eth0
? (192.168.5.207) at 00:50:56:b0:80:64 [ether] on eth0
? (192.168.5.245) at 00:50:56:b0:e2:13 [ether] on eth0
? (192.168.5.43) at 00:50:56:b0:0f:52 [ether] on eth0
? (192.168.5.54) at 00:50:56:b0:e4:6d [ether] on eth0
? (192.168.5.93) at 00:50:56:b0:1b:5b [ether] on eth0
```

Thus `arp -n` gives you the simplest quick overview of the veth pairs associated with the pods on your local machine.

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

In both clusters, run `kubectl scale deployment coredns --replicas=10 -n kube-system`.  Then re run these commands.  You'll see new ip a entries for the containers.

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

This command works the same in antrea.  We can see that there are lots of packets reaching the bridge, compared with the other containers: 
- 2,801,928 received packets for the `antrea_gw0` device (everything goes through ovs)
- 831670 received packets for the `geneve_sys` device
- 383670 received packets for an 1 coredns container
- 384203 received packets for the other coredns container

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
```
 +-----------------------------------+
 |                                   |    +--------+
 +----------+                        |    |        |
 |          |                    +------->+ 172...3|
 |	192.168  |   +-----------|   +    +--------+
 | .9.130   +-------->   tun0   +    |
 |          |        +----------+    |    +--------+
 +----------+        |   tun0   |    |    | 172...5|
 |                   +----------+    |    +--------+
 |                   |   tun0   |    |
 |                   +----------+    |
 |                   |   tun0   |    |    +--------+
 |                   +----------+    |    | 172...4|
 |                                   |    |        |
 +------------+                      |    +--------+
 || 192.168   |       +-----------+  |
 || .173.64   +------>+ locally   |  |    +--------+
 ||           |       | connected |  |    | 172...2|
 +------------+       | (0.0.0.0) |  |    |        |
 |                    +-----------+  |    +--------+
 |                                   |
 |                                   |
++-----------------------------------+

 +-----------------------------------+
 |              +-------------+      |    +--------+
 |              |   antrea|gw0+---------> |  ovs   |
 | 100.96.      +-------------+      |    |gateway |
 | 0.2   +----->+   antrea|gw0|      |    |  0.1   |
 |              +-------------+      |    +--------+
 |                                   |
 |                                   |
 |                                   |
 |                                   |
 +-----------------------------------+
```

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

We've looked at:

- How the host maps routes Pod traffic
- How, from the host, you can verify incoming pod traffic

However, we havent yet looked at *services*.

IPTables (or IPvS) ultimately give us the ability to see how traffic is being routed from services to pods.

### Making sure service routing is working

The simplest thing you can do to start is look for all service endpoints:

`iptables-save`

From here, you can look for the comment rules, which will tell you the services which are associated with a rule.

```
-A KUBE-SVC-TCOU7JCQXEZGVUNU -m comment --comment "kube-system/kube-dns:dns" -m statistic --mode random --probability 0.10000000009 -j KUBE-SEP-QIVPDYSUOLOYQCAA
-A KUBE-SVC-TCOU7JCQXEZGVUNU -m comment --comment "kube-system/kube-dns:dns" -m statistic --mode random --probability 0.11111111101 -j KUBE-SEP-N76EJY3A4RTXTN2I
-A KUBE-SVC-TCOU7JCQXEZGVUNU -m comment --comment "kube-system/kube-dns:dns" -m statistic --mode random --probability 0.12500000000 -j KUBE-SEP-LSGM2AJGRPG672RM
```

After looking at these services, you'll want to find the corresponding `SEP` rules for them:

```
root@calico-worker3:/# iptables-save | grep SEP-QI
:KUBE-SEP-QIVPDYSUOLOYQCAA - [0:0]
### Masquerading happens here for outgoing traffic...
-A KUBE-SEP-QIVPDYSUOLOYQCAA -s 192.168.143.65/32 -m comment --comment "kube-system/kube-dns:dns" -j KUBE-MARK-MASQ
-A KUBE-SEP-QIVPDYSUOLOYQCAA -p udp -m comment --comment "kube-system/kube-dns:dns" -m udp -j DNAT --to-destination 192.168.143.65:53
-A KUBE-SVC-TCOU7JCQXEZGVUNU -m comment --comment "kube-system/kube-dns:dns" -m statistic --mode random --probability 0.10000000009 -j KUBE-SEP-QIVPDYSUOLOYQCAA
```

This step is *the exact same* in any CNI provider you use, so we dont provide an antrea/calico comparison.

### Looking at network policies

To start looking at how network policies might be affecting traffic, run a networkpolicy test:

```
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: web-deny-all
spec:
  podSelector:
    matchLabels:
      app: web
  ingress: []
```

A good way to uniformly test these policies is to create a daemonset running the same container in all nodes:

```
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: nginx-ds
spec:
  selector:
    matchLabels:
      app: web
  template:
    metadata:
      labels:
        app: web
    spec:
      containers:
      - name: nginx
        image: nginx
```


#### How are these policies implemented? 

NetworkPolicies are often but not always implemented by your CNI providers as iptables rules.

In calico, you'll see policies such as this.  This is where the 'drop' rule of a policy is implemented.
If you were to enable a policy later, these rules would get superceded by those add rules, per the NetworkPolicy API specification.

```
> -A cali-tw-calic5cc839365a -m comment --comment "cali:Uv2zkaIvaVnFWYI9" -m comment --comment "Start of policies" -j MARK --set-xmark 0x0/0x20000
> -A cali-tw-calic5cc839365a -m comment --comment "cali:7OLyCb9i6s_CPjbu" -m mark --mark 0x0/0x20000 -j cali-pi-_IDb4Gbl3P1MtRtVzfEP
> -A cali-tw-calic5cc839365a -m comment --comment "cali:DBkU9PXyu2eCwkJC" -m comment --comment "Return if policy accepted" -m mark --mark 0x10000/0x10000 -j RETURN
> -A cali-tw-calic5cc839365a -m comment --comment "cali:tioNk8N7f4P5Pzf4" -m comment --comment "Drop if no policies passed packet" -m mark --mark 0x0/0x20000 -j DROP
> -A cali-tw-calic5cc839365a -m comment --comment "cali:wcGG1iiHvTXsj5lq" -j cali-pri-kns.default
> -A cali-tw-calic5cc839365a -m comment --comment "cali:gaGDuGQkGckLPa4H" -m comment --comment "Return if profile accepted" -m mark --mark 0x10000/0x10000 -j RETURN
> -A cali-tw-calic5cc839365a -m comment --comment "cali:B6l_lueEhRWiWwnn" -j cali-pri-ksa.default.default
> -A cali-tw-calic5cc839365a -m comment --comment "cali:McPS2ZHiShhYyFnW" -m comment --comment "Return if profile accepted" -m mark --mark 0x10000/0x10000 -j RETURN
> -A cali-tw-calic5cc839365a -m comment --comment "cali:lThI2kHuPODjvF4v" -m comment --comment "Drop if no profiles matched" -j DROP
```

Antrea also implements networkpolicies, but uses openvswitch flows, and writes these flows to table 90. 

In antrea, running a similar workload, you'll see these policies created.  An easy way to do this is to call ovs-ofctl .  Typically
this is done from inside a container since antrea-agents are fully configured with all openvswitch administrative binaries.
This can also work from the host as well if needed, just install the ovs utilities as needed... 

```
/tmp » kubectl -n kube-system exec -it antrea-agent-2kksz ovs-ofctl dump-flows br-int | grep table=90                 
```

```
kubectl exec [POD] [COMMAND] is DEPRECATED and will be removed in a future version. Use kubectl kubectl exec [POD] -- [COMMAND] instead.
Defaulting container name to antrea-agent.
Use 'kubectl describe pod/antrea-agent-2kksz -n kube-system' to see all of the containers in this pod.
 cookie=0x2000000000000, duration=344936.777s, table=90, n_packets=0, n_bytes=0, priority=210,ct_state=-new+est,ip actions=resubmit(,105)
 cookie=0x2000000000000, duration=344936.776s, table=90, n_packets=83160, n_bytes=6153840, priority=210,ip,nw_src=100.96.26.1 actions=resubmit(,105)
 
 ### This line shows that you have some pods which are being allowed, via the ovs flow register, into the cluster.... 
 cookie=0x2050000000000, duration=22.296s, table=90, n_packets=0, n_bytes=0, priority=200,ip,reg1=0x18 actions=conjunction(1,2/2)
 
 cookie=0x2050000000000, duration=22.300s, table=90, n_packets=0, n_bytes=0, priority=190,conj_id=1,ip actions=load:0x1->NXM_NX_REG6[],resubmit(,105)
 cookie=0x2000000000000, duration=344936.782s, table=90, n_packets=149662, n_bytes=11075281, priority=0 actions=resubmit(,100)
```









































# TODO

- Calico Antrea cni swap overlay 
- ip route trace 

