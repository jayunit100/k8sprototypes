# Tanzu Framework 

This package gives us some basic "Tanzu Framework" primitives we can work with 
when scaffolding ginkgo tests.  In general we expect the tests to:

- Create infrastructure, somehow, using the Testbeds package
- Create a few model classes, using the model classes, which represent tanzu cluster elements
- Use this package, tanzu framework, to take those model classes (as inputs), and then run those 
commands (either in a simulated testbed, or in a real testbed) in order to create a fully working cluster.
