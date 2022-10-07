snippets reusable for other python scripts.
used to make https://github.com/jayunit100/k8sprototypes/blob/master/sonoeasy/sonoeasy.py, which
parses sonobuoy records and makes text summaries like this:

```

                shrt     avg     dev      long   cnt     totaltime(s)
sig-api-machin   0.2s    9.1     10.2    63.4s   63      571.0
sig-instrum      0.2s    0.2     0.0     0.2s    4       0.9
sig-netw         0.2s    9.4     9.8     35.5s   41      385.1
sig-scheduling   1.3s    69.1    95.1    304.5s  9       622.0
sig-auth         0.2s    6.4     12.4    34.3s   7       44.8
sig-apps         0.4s    30.6    64.7    326.2s  47      1439.6
sig-cli          0.2s    3.8     3.7     13.7s   17      65.3
sig-storage      0.2s    7.1     11.5    82.9s   81      574.0
sig-node         0.2s    22.9    53.9    243.5s  77      1760.1

```
