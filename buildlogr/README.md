#  How to customize

To get the get-prow.sh script working, look at a large list of test runs, like https://testgrid.k8s.io/sig-network-gce#presubmit-network-policies,%20google-gce .

Then just copy the URLs from the raw-build.log txt.

# directions


1) Modify prow-get.sh to grab the build log files your interested in
2) Then run the ./clusters.sh script.  It will scrape down those files and find failures that happened near the 'failed' tests, temporily

The output is a list of frequency sorted 'suspects' which might be correlated to failures.  These might be tests which

- use alot of resources
- take a long time to run
- interfer with api calls or add security interferences
- serailzed the cluster performance in some other way


