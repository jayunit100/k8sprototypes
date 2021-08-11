# NOTE https://gist.github.com/aojea/097b5a8418fbbcb2b55e72a4cf6e62f7 has a really good general
# script to be reused, we should adopt that in the future

#!/bin/bash
sha=`curl -L http://dl.k8s.io/ci/latest.txt`

# TODO ~ add a way to remeber sha from older 1.22 or 1.21 e2es... 
# wget  https://storage.googleapis.com/kubernetes-release-dev/ci/${sha}/kubernetes-test-linux-amd64.tar.gz
wget https://storage.googleapis.com/k8s-release-dev/ci/${sha}/kubernetes-test-linux-amd64.tar.gz
tar -xvf kubernetes-test-linux-amd64.tar.gz -C /tmp/
cp /tmp/kubernetes/test/bin/e2e.test ./e2e.test
chmod 755 ./e2e.test
echo "e2e.test is ready to run now at `pwd`/e2e.test"
./e2e.test --version
#./e2e.test --provider=local --kubeconfig=/home/kubo/.kube/config \
#--dump-logs-on-failure=false \
#--ginkgo.focus='should be able to create a functioning NodePort service for Windows' \
#--ginkgo.skip='Slow' \
#--node-os-distro=windows
