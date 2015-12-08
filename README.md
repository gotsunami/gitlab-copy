
# GitLab Copy

`gitlab-copy` is a simple tool for copying issues/labels/milestones/notes from one GitLab project to another, possibly running on different GitLab instances.

Grab a binary version for your platform and write a YAML config file to specify source and target projects:

```
$ ./gitlab-copy -h
Usage: ./bin/gitlab-copy [options] configfile

Where configfile holds YAML data like:
from:
    url: https://gitlab.mydomain.com
    token: atoken
    project: namespace/project
to:
    url: https://gitlab.myotherdomain.com
    token: anothertoken
    project: namespace/project

Options:
  -version
        
  -y    apply migration for real
```

## Compile From Source

Install `gb` to compile the project:
```
$ go get github.com/constabulary/gb/...
```

Then build the project:
```
$ gb build
```

The program gets compiled into `bin/gitlab-copy`.
