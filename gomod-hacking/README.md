# Go mod example, minimal

The purpose of this is to do some experiments on how go.mod works.

interestingly, the go.mod file generated for the webapp directory

has version generation from time.Time, but `cobra` which lives on the internet, 
has a versino generated from 1.2.1.

Asked on stackoverflow: How does go infer the proper `go mod tidy` version for a local repo?


```
require (
        example.com/db_api v0.0.0-00010101000000-000000000000
        github.com/spf13/cobra v1.2.1
)
```
