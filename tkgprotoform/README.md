# Welcome to TKG PRotoform, a self documenting CI and TKG installation platform

```
go build cmd/main/main.go  
```

Then 

```
mkdir geetika/
cp main geetika/tkg-protoform 
./tkg-protoform
```

Youll see a bunch of TKG custom installation media (configurable via viper)

```
➜  tkgprotoform git:(f6f154e) ✗ tree /tmp/geetika-testbed
/tmp/geetika-testbed
├── corgon_cluster.yaml
└── management_cluster.sh
...
```

Now run the scripts, and TKG will be up and running.  

The "secret sauce" is the use of the `//go:embed management-cluster.sh` directive.
This directive takes the prototypical scripts, parses and substitutes them (that
logic is in the works), and then regurgitates them.

Thus, its a perfectly reproducible and portable, "hermetically sealed" in a single binary,
installation and testing tool for TKG and tanzu framework, including the entire end user experience.

# Why not write installers in a programming language? 

Because we want to autogenerate our documentation from this platform.   Otherwise, we'll have

- one set of automation that lives in golang, python, groovy, junit whatever
- another set of psuedo automation that lives in some mystical combination of READMEs, and internet docs, and...
- scott rosenberg blog posts.

