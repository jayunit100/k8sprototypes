# Performance test

- iperf test
  - image : k8s.gcr.io/e2e-test-images/agnhost:2.21
  - iperf server : iperf -s -p 9999 -D
  - iperf client : iperf -p 9999 -c server-ip

- Create/Delete pod 100 times:
  - ./run_deployment.sh 100

- Create service:
  - ./run_servcie.sh nodeport
  - ./run_service.sh clusterip

- Create network policy
  - create_net_policies.py > net_policies.yaml
  - kubectl apply -f test-pod1.yaml
  - kubectl apply -f test-pod2.yaml
  - kubectl apply -f net_policies.yaml
