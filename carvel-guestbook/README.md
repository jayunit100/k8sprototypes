The directories in here each iteratively
build up the "guestbook" application in a 
simple way.

- v0/ is just the yaml
- v1/ is the separated yaml files which can be combined by `ytt -f` to make a composite guestbook.yaml file
- v2/ is an example of how to overlay onto the output of v1/ using ytt
- v3/ is an example of how to integrate kapp into the mix so that the objects are managed as a single application by the `kapp` tool

Just an experiment in how to iteratively introduce carvel into your life in a sustainable way :) 

