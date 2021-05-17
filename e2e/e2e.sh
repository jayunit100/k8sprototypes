#!/bin/bash
sha=`curl -L http://dl.k8s.io/ci/latest.txt`
wget  https://storage.googleapis.com/kubernetes-release-dev/ci/${sha}/kubernetes-test-linux-amd64.tar.gz
tar -xvf kubernetes-test-linux-amd64.tar.gz -C /tmp/
cp /tmp/e2e.test ./e2e.test
chmod 755 ./e2e.test
echo "e2e.test is ready to run now"
./e2e.test --version
#./e2e.test --provider=local --kubeconfig=/home/kubo/.kube/config \
#--dump-logs-on-failure=false \
#--ginkgo.focus='should be able to create a functioning NodePort service for Windows' \
#--ginkgo.skip='Slow' \
#--node-os-distro=windows
