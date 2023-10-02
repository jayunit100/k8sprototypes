### Install e2e binary

wget https://github.com/jayunit100/k8sprototypes/raw/master/e2e/e2e.sh 
chmod 777 ./e2e.sh
./e2e.sh

### Run iperf in parallel

run_test() {                                                                                                                                               │
  while true; do                                                                                                                                           │
    ./e2e.test --provider=local --kubeconfig=/home/kubo/.kube/config --ginkgo.focus="iperf*" --ginkgo.skip="Disruptive" -v=5 --dump-logs-on-failure=false -│
-allowed-not-ready-nodes=5 --minStartupPods=-1                                                                                                             │
  done                                                                                                                                                     │
}                                                                                                                                                          │
                                                                                                                                                           │
# Run 10 instances in parallel                                                                                                                             │
for i in {1..10}; do                                                                                                                                       │
  run_test &                                                                                                                                               │
done                                                                                                                                                       │
                                                                                                                                                           │
# Wait for all background processes to finish                                                                                                              │
wait           
