---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-0
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod1
  policyTypes:
  - Ingress
  ingress:
  - from:
    ports:
    - protocol: TCP
      port: 10000
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-1
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod1
  policyTypes:
  - Ingress
  ingress:
  - from:
    ports:
    - protocol: TCP
      port: 10020
      endPort: 10039
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-2
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod1
  policyTypes:
  - Ingress
  ingress:
  - from:
    ports:
    - protocol: UDP
      port: 10040
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-3
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod1
  policyTypes:
  - Ingress
  ingress:
  - from:
    ports:
    - protocol: UDP
      port: 10060
      endPort: 10079
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-4
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod1
  policyTypes:
  - Ingress
  ingress:
  - to:
    ports:
    - protocol: TCP
      port: 10080
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-5
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod1
  policyTypes:
  - Ingress
  ingress:
  - to:
    ports:
    - protocol: TCP
      port: 10100
      endPort: 10119
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-6
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod1
  policyTypes:
  - Ingress
  ingress:
  - to:
    ports:
    - protocol: UDP
      port: 10120
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-7
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod1
  policyTypes:
  - Ingress
  ingress:
  - to:
    ports:
    - protocol: UDP
      port: 10140
      endPort: 10159
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-8
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod1
  policyTypes:
  - Egress
  egress:
  - from:
    ports:
    - protocol: TCP
      port: 10160
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-9
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod1
  policyTypes:
  - Egress
  egress:
  - from:
    ports:
    - protocol: TCP
      port: 10180
      endPort: 10199
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-10
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod1
  policyTypes:
  - Egress
  egress:
  - from:
    ports:
    - protocol: UDP
      port: 10200
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-11
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod1
  policyTypes:
  - Egress
  egress:
  - from:
    ports:
    - protocol: UDP
      port: 10220
      endPort: 10239
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-12
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod1
  policyTypes:
  - Egress
  egress:
  - to:
    ports:
    - protocol: TCP
      port: 10240
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-13
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod1
  policyTypes:
  - Egress
  egress:
  - to:
    ports:
    - protocol: TCP
      port: 10260
      endPort: 10279
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-14
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod1
  policyTypes:
  - Egress
  egress:
  - to:
    ports:
    - protocol: UDP
      port: 10280
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-15
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod1
  policyTypes:
  - Egress
  egress:
  - to:
    ports:
    - protocol: UDP
      port: 10300
      endPort: 10319
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-16
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod2
  policyTypes:
  - Ingress
  ingress:
  - from:
    ports:
    - protocol: TCP
      port: 10320
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-17
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod2
  policyTypes:
  - Ingress
  ingress:
  - from:
    ports:
    - protocol: TCP
      port: 10340
      endPort: 10359
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-18
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod2
  policyTypes:
  - Ingress
  ingress:
  - from:
    ports:
    - protocol: UDP
      port: 10360
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-19
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod2
  policyTypes:
  - Ingress
  ingress:
  - from:
    ports:
    - protocol: UDP
      port: 10380
      endPort: 10399
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-20
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod2
  policyTypes:
  - Ingress
  ingress:
  - to:
    ports:
    - protocol: TCP
      port: 10400
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-21
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod2
  policyTypes:
  - Ingress
  ingress:
  - to:
    ports:
    - protocol: TCP
      port: 10420
      endPort: 10439
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-22
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod2
  policyTypes:
  - Ingress
  ingress:
  - to:
    ports:
    - protocol: UDP
      port: 10440
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-23
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod2
  policyTypes:
  - Ingress
  ingress:
  - to:
    ports:
    - protocol: UDP
      port: 10460
      endPort: 10479
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-24
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod2
  policyTypes:
  - Egress
  egress:
  - from:
    ports:
    - protocol: TCP
      port: 10480
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-25
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod2
  policyTypes:
  - Egress
  egress:
  - from:
    ports:
    - protocol: TCP
      port: 10500
      endPort: 10519
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-26
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod2
  policyTypes:
  - Egress
  egress:
  - from:
    ports:
    - protocol: UDP
      port: 10520
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-27
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod2
  policyTypes:
  - Egress
  egress:
  - from:
    ports:
    - protocol: UDP
      port: 10540
      endPort: 10559
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-28
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod2
  policyTypes:
  - Egress
  egress:
  - to:
    ports:
    - protocol: TCP
      port: 10560
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-29
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod2
  policyTypes:
  - Egress
  egress:
  - to:
    ports:
    - protocol: TCP
      port: 10580
      endPort: 10599
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-30
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod2
  policyTypes:
  - Egress
  egress:
  - to:
    ports:
    - protocol: UDP
      port: 10600
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-31
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod2
  policyTypes:
  - Egress
  egress:
  - to:
    ports:
    - protocol: UDP
      port: 10620
      endPort: 10639
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-32
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod1
  policyTypes:
  - Ingress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          role: test-pod2
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-33
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod1
  policyTypes:
  - Egress
  egress:
  - to:
    - podSelector:
        matchLabels:
          role: test-pod2
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-34
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod2
  policyTypes:
  - Ingress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          role: test-pod1
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-35
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod2
  policyTypes:
  - Egress
  egress:
  - to:
    - podSelector:
        matchLabels:
          role: test-pod1
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-36
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod1
  policyTypes:
  - Ingress
  ingress:
  - from:
    - ipBlock:
        cidr: 100.36.0.0/16
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-37
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod1
  policyTypes:
  - Ingress
  ingress:
  - from:
    - ipBlock:
        cidr: 100.37.0.0/16
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-38
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod1
  policyTypes:
  - Egress
  egress:
  - to:
    - ipBlock:
        cidr: 100.38.0.0/16
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-39
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod1
  policyTypes:
  - Egress
  egress:
  - to:
    - ipBlock:
        cidr: 100.39.0.0/16
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-40
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod2
  policyTypes:
  - Ingress
  ingress:
  - from:
    - ipBlock:
        cidr: 100.0.0.0/8
        except:
        - 100.40.0.0/16
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-41
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod2
  policyTypes:
  - Egress
  egress:
  - to:
    - ipBlock:
        cidr: 100.0.0.0/8
        except:
        - 100.41.0.0/16
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-42
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod1
  policyTypes:
  - Ingress
  ingress:
  - from:
    - ipBlock:
        cidr: 100.42.0.0/16
    ports:
    - protocol: TCP
      port: 10840
      endPort: 10859
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-43
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod1
  policyTypes:
  - Ingress
  ingress:
  - from:
    - ipBlock:
        cidr: 100.43.0.0/16
    ports:
    - protocol: UDP
      port: 10860
      endPort: 10879
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-44
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod1
  policyTypes:
  - Egress
  egress:
  - to:
    - ipBlock:
        cidr: 100.44.0.0/16
    ports:
    - protocol: TCP
      port: 10880
      endPort: 10899
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-45
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod1
  policyTypes:
  - Egress
  egress:
  - to:
    - ipBlock:
        cidr: 100.45.0.0/16
    ports:
    - protocol: UDP
      port: 10900
      endPort: 10919
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-46
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod1
  policyTypes:
  - Ingress
  ingress:
  - from:
    - ipBlock:
        cidr: 100.46.0.0/16
    - ipBlock:
        cidr: 100.47.0.0/16
    ports:
    - protocol: TCP
      port: 10920
      endPort: 10939
    - protocol: TCP
      port: 10940
      endPort: 10959
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-47
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod1
  policyTypes:
  - Ingress
  ingress:
  - from:
    - ipBlock:
        cidr: 100.48.0.0/16
    - ipBlock:
        cidr: 100.49.0.0/16
    ports:
    - protocol: UDP
      port: 10960
      endPort: 10979
    - protocol: UDP
      port: 10980
      endPort: 10999
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-48
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod1
  policyTypes:
  - Egress
  egress:
  - to:
    - ipBlock:
        cidr: 100.50.0.0/16
    - ipBlock:
        cidr: 100.51.0.0/16
    ports:
    - protocol: TCP
      port: 11000
      endPort: 11019
    - protocol: TCP
      port: 11020
      endPort: 11039
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-49
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod1
  policyTypes:
  - Egress
  egress:
  - to:
    - ipBlock:
        cidr: 100.52.0.0/16
    - ipBlock:
        cidr: 100.53.0.0/16
    ports:
    - protocol: UDP
      port: 11040
      endPort: 11059
    - protocol: UDP
      port: 11060
      endPort: 11079
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-50
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod1
  policyTypes:
  - Ingress
  ingress:
  - from:
    - ipBlock:
        cidr: 100.54.0.0/16
    - ipBlock:
        cidr: 100.55.0.0/16
    ports:
    - protocol: TCP
      port: 11080
      endPort: 11099
    - protocol: UDP
      port: 11100
      endPort: 11119
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-netpolicy-51
  namespace: default
spec:
  podSelector:
    matchLabels:
      role: test-pod1
  policyTypes:
  - Egress
  egress:
  - to:
    - ipBlock:
        cidr: 100.56.0.0/16
    - ipBlock:
        cidr: 100.57.0.0/16
    ports:
    - protocol: TCP
      port: 11120
      endPort: 11139
    - protocol: UDP
      port: 11140
      endPort: 11159
