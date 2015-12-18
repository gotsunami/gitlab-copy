
# GitLab Copy

`gitlab-copy` is a simple tool for copying issues/labels/milestones/notes from one GitLab project to another, possibly running on different GitLab instances.

`gitlab-copy` won't copy anything until told **explicitely** to do so. Running it from the command line will show some stats only.

[Grab a binary version](https://github.com/gotsunami/gitlab-copy/releases/latest) for your platform or have a look to the **Compile From Source** section.

## Features

Please note this is **beta** software. The following features are available:

- Creates milestones if not existing on target
- Creates labels if not existing on target
- Creates issues if not existing on target (by title)
- Apply closed status on issues, if any
- Creates notes (attached to issues)
- Can specify in the config file a specific issue or range of issues to copy

## Usage

First, write a `gitlab.yml` YAML config file to specify source and target projects, along with your GitLab account tokens:
```
from:
    url: https://gitlab.mydomain.com
    token: atoken
    project: namespace/project
to:
    url: https://gitlab.myotherdomain.com
    token: anothertoken
    project: namespace/project
```

Note that a specific issue or ranges of issues can be specified in the YAML config file. If you want to
copy only issue #15 and issues #20 to #30, add an `issues` key in the `from:` key:
```
from:
    url: https://gitlab.mydomain.com
    token: atoken
    project: namespace/project
    issues:
        - 15
        - 20-30
...
```

Now grab some project stats by running
```
$ ./gitlab-copy gitlab.yml
```

If everything looks good, run the same command, this time with the `-y` flag to effectively copie issues between GitLab
instances (they can be the same):
```
$ ./gitlab-copy -y gitlab.yml
```

## Compile From Source

Install `gb` to compile the project:
```
$ go get github.com/constabulary/gb/...
```

- To build the project, just run `make`. The program gets compiled into `bin/gitlab-copy`
- Cross-compile with `make buildall`
- Prepare distribution packages with `make dist`

## License

MIT. See `LICENSE` file.
