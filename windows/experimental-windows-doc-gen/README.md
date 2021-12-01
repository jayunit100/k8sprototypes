# Experimental (WIP) Generator for installation docs

This is an experiment around what it would be like to use ytt to GENERATE
documentation, as well as CI systems that are self testable, for VMWare tanzu.
I might have time to get the whole thing working, but for now heres how it works:


## Generate a CI runnable script that installs tanzu

```
ytt -f defaults.yaml -f instructions.yaml \
grep RUNME \
cut -d'-' -f2 \
sed "s/'//g"
```
The above:
1) Runs ytt to make yaml structured, readable docs
2) Greps out the "RUNnable" commands
3) removes the dashes (from the yaml)
4) removes the single qoutes (from the yaml)

Thus, it outputs an executable script :)

## Generate the raw yaml for the above script, to generate docs about installing tanzu

Will generate an executable "script" which a CI system could use, but without the grep:

```
ytt -f defaults.yaml -f instructions.yaml
```

You essentially get product documentation that can be copied into formal documentation. 

output (docs)
```
steps:
- step: Install Mgmt cluster
  steps:
  - download tarball:
    - wget www.tanzu.com/files/tanzu.tar.gz
  - unzip tarball:
    - tar -xvf tanzu.tar.gz
  - define vsphere credentials:
    - export GOVC_VSPHERE_USER=administrator@vsphere.local
    - export GOVC_VSPHERE_PASSWORD=Admin!23
- run tanzu init
- step: Install Windows cluster
  steps:
  - make credentials.json
  - make windows.json
  - update autounattend if needed
  - craft docker command w/ subs
```

output (ci)
```
 'wget www.tanzu.com/files/tanzu.tar.gz # RUNME'
 'tar
 'export GOVC_VSPHERE_USER=administrator@vsphere.local # RUNME '
 'export GOVC_VSPHERE_PASSWORD=Admin!23 # RUNME '
 'make credentials.json # RUNME'
 'make windows.json # RUNME'
 'update autounattend if needed # RUNME'
 'craft docker command w/ subs # RUNME'
```


