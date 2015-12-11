
# GitLab Copy

`gitlab-copy` is a simple tool for copying issues/labels/milestones/notes from one GitLab project to another, possibly running on different GitLab instances.

`gitlab-copy` won't copy anything until told **explicitely** to do so. Running it from the command line will show some stats only.

[Grab a binary version](https://github.com/gotsunami/gitlab-copy/releases/latest) for your platform and write a YAML config file to specify source and target projects:

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

## Examples

*TBW*

## Compile From Source

Install `gb` to compile the project:
```
$ go get github.com/constabulary/gb/...
```

- To build the project, just run `make`. The program gets compiled into `bin/gitlab-copy`
- Cross-compile with `make buildall`
- Prepare distribution packages with `make dist`
