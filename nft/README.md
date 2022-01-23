# NFT !!!

To spin up an NFT environment to play with... run...
`vagrant up`

## IP addresses

get the IP addresses of the nodes:

`ip a | grep 192`

- 192.168.70.7
- 192.168.70.8

## Setup 

Run `sudo apt get install nftables`

Then you can run `nft --help`

## tables

- family
  - table
    - chains
      - *DECLARING* them: 
        - type: 
          - filter works w/ tables=(ip, ip6 and inet) + (arp bridge)
  	  - route works w/ tables=(ip, ip6 and inet)
	  - nat works w/ tables=(ip, ip6)
        - hooks for BASE CHAIN:
	  - ingress: before pre routing, right after the NIC
	  - prerouting: ALL incoming packets, after INGRESS
	  - input: addressed to the LOCAL system 
	  - fwd: packets NOT the local system 
	  - output: local packets going OUT of the machine
	  - postrouting: all packets RIGHT BEFORE THEY leave the machine
      - *ADDING* RULES to them:
        
    - sets       
    - maps
    - flowtables
    - stateful objects


