#!/usr/bin/python3

# {ingress, egress}
# {allow, deny}
# {tcp, udp, ip, podSelector, namespaceSelector}
# {to, from}
# {w/, w/o} 

import os

def create_policy_header(name, pod, ingress):
    ingr = "ingress" if ingress else "egress"
    pt = "Ingress" if ingress else "Egress"
    hdr = f"""apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: {name}
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: {pod}
  policyTypes:
  - {pt}
  {ingr}:"""
    return hdr

def create_allow_port(name, pod, ingress, fromto, proto, port):
    hdr = create_policy_header(name, pod, ingress)
    s = f"""
  - {fromto}:
    ports:
    - protocol: {proto}
      port: {port}"""
    return hdr + s

def create_allow_ports(name, pod, ingress, fromto, proto, port, endPort):
    hdr = create_policy_header(name, pod, ingress)
    s = f"""
  - {fromto}:
    ports:
    - protocol: {proto}
      port: {port}
      endPort: {endPort}"""
    return hdr + s

def create_allow_pod(name, pod, ingress, fromto, pod2):
    hdr = create_policy_header(name, pod, ingress)
    s = f"""
  - {fromto}:
    - podSelector:
        matchLabels:
          role: {pod2}"""
    return hdr + s

def create_allow_ips(name, pod, ingress, fromto, ipBlock):
    hdr = create_policy_header(name, pod, ingress)
    s = f"""
  - {fromto}:
    - ipBlock:
        cidr: {ipBlock}"""
    return hdr + s

def create_deny_ips(name, pod, ingress, fromto, ipBlock):
    hdr = create_policy_header(name, pod, ingress)
    s = f"""
  - {fromto}:
    - ipBlock:
        cidr: 100.0.0.0/8
        except:
        - {ipBlock}"""
    return hdr + s

def create_allow_ips_ports_proto(name, pod, ingress, ipBlock, port, endPort, proto):
    hdr = create_policy_header(name, pod, ingress)
    fromto = "from" if ingress else "to"
    s = f"""
  - {fromto}:
    - ipBlock:
        cidr: {ipBlock}
    ports:
    - protocol: {proto}
      port: {port}
      endPort: {endPort}"""
    return hdr + s

def create_allow_ips_ports_proto2(name, pod, ingress, ipBlocks, ports, endPorts, proto):
    hdr = create_policy_header(name, pod, ingress)
    fromto = "from" if ingress else "to"
    s = f"""
  - {fromto}:
    - ipBlock:
        cidr: {ipBlocks[0]}
    - ipBlock:
        cidr: {ipBlocks[1]}
    ports:
    - protocol: {proto}
      port: {ports[0]}
      endPort: {endPorts[0]}
    - protocol: {proto}
      port: {ports[1]}
      endPort: {endPorts[1]}"""
    return hdr + s

def create_allow_ips_ports(name, pod, ingress, ipBlocks, ports, endPorts):
    hdr = create_policy_header(name, pod, ingress)
    fromto = "from" if ingress else "to"
    s = f"""
  - {fromto}:
    - ipBlock:
        cidr: {ipBlocks[0]}
    - ipBlock:
        cidr: {ipBlocks[1]}
    ports:
    - protocol: TCP
      port: {ports[0]}
      endPort: {endPorts[0]}
    - protocol: UDP
      port: {ports[1]}
      endPort: {endPorts[1]}"""
    return hdr + s

def get_port_range(idx):
    port = 10000 + idx * 20
    endPort = port + 20 - 1
    return port, endPort

def get_ip_block(idx):
    return f"100.{idx}.0.0/16"

def print_policy(s):
    print("---")
    print(s)

idx = 0
pods = ["test-pod1", "test-pod2"]
netpolicy_prefix = "test-netpolicy-"

for pod in pods:
    for ingr in [True, False]:
        for ft in ["from", "to"]:
            for proto in ["TCP", "UDP"]:
                port, endPort = get_port_range(idx)
                s = create_allow_port(f"{netpolicy_prefix}{idx}", pod, ingr, ft, proto, port)
                idx = idx + 1
                print_policy(s)
                port, endPort = get_port_range(idx)
                s = create_allow_ports(f"{netpolicy_prefix}{idx}", pod, ingr, ft, proto, port, endPort)
                idx = idx + 1
                print_policy(s)

for podidx in [0, 1]:
    for ingr in [True, False]:
        ft = "from" if ingr else "to"
        s = create_allow_pod(f"{netpolicy_prefix}{idx}", pods[podidx], ingr, ft, pods[1-podidx])
        idx = idx + 1
        print_policy(s)

for ingr in [True, False]:
    ft = "from" if ingr else "to"
    for r in [0, 1]:
        ipBlock = get_ip_block(idx)
        s = create_allow_ips(f"{netpolicy_prefix}{idx}", pods[0], ingr, ft, ipBlock)
        idx = idx + 1
        print_policy(s)

for ingr in [True, False]:
    ft = "from" if ingr else "to"
    ipBlock = get_ip_block(idx)
    s = create_deny_ips(f"{netpolicy_prefix}{idx}", pods[1], ingr, ft, ipBlock)
    idx = idx + 1
    print_policy(s)

for ingr in [True, False]:
    for proto in ["TCP", "UDP"]:
        ipBlock = get_ip_block(idx)
        port, endPort = get_port_range(idx)
        s = create_allow_ips_ports_proto(f"{netpolicy_prefix}{idx}", pods[0], ingr, ipBlock, port, endPort, proto)
        idx = idx + 1
        print_policy(s)

ip_idx = idx
port_idx = idx

for ingr in [True, False]:
    for proto in ["TCP", "UDP"]:
        ipBlock0 = get_ip_block(ip_idx)
        ipBlock1 = get_ip_block(ip_idx+1)
        port0, endPort0 = get_port_range(port_idx)
        port1, endPort1 = get_port_range(port_idx+1)
        ipBlocks = [ipBlock0, ipBlock1]
        ports = [port0, port1]
        endPorts = [endPort0, endPort1]
        s = create_allow_ips_ports_proto2(f"{netpolicy_prefix}{idx}", pods[0], ingr, ipBlocks, ports, endPorts, proto)
        idx = idx + 1
        ip_idx = ip_idx + 2
        port_idx = port_idx + 2
        print_policy(s)

for ingr in [True, False]:
    ipBlock0 = get_ip_block(ip_idx)
    ipBlock1 = get_ip_block(ip_idx+1)
    port0, endPort0 = get_port_range(port_idx)
    port1, endPort1 = get_port_range(port_idx+1)
    ipBlocks = [ipBlock0, ipBlock1]
    ports = [port0, port1]
    endPorts = [endPort0, endPort1]
    s = create_allow_ips_ports(f"{netpolicy_prefix}{idx}", pods[0], ingr, ipBlocks, ports, endPorts)
    idx = idx + 1
    ip_idx = ip_idx + 2
    port_idx = port_idx + 2
    print_policy(s)
