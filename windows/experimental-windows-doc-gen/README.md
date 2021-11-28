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
