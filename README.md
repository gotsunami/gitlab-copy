
# GitLab Copy

[![Build Status](https://travis-ci.org/gotsunami/gitlab-copy.svg?branch=master)](https://travis-ci.org/gotsunami/gitlab-copy)

`gitlab-copy` is a simple tool for copying issues/labels/milestones/notes from one GitLab project to another, possibly running on different GitLab instances.

`gitlab-copy` won't copy anything until told **explicitely** to do so. Running it from the command line will show some stats only.

**Note**: GitLab 8.6 introduced the ability [to move an issue to another project](https://about.gitlab.com/2016/03/22/gitlab-8-6-released/), but on the same GitLab installation only. `gitlab-copy` can still prove valuable to move issues between projects on different GitLab hosts and to perform batch operations from the command line (see the feature list below).

## Download

Installing `gitlab-copy` is very easy since it comes as a static binary with no dependencies. Just [grab a compiled version](https://github.com/gotsunami/gitlab-copy/releases/latest) for your platform (or have a look to the **Compile From Source** section).

## Features

The following features are available:

- Copy milestones if not existing on target
- Copy all source labels on target (use `labelsOnly` to copy labels only, see below)
- Copy issues if not existing on target (by title)
- Apply closed status on issues, if any
- Set issue's assignee (if user exists) and milestone, if any
- Copy notes (attached to issues), preserving user ownership
- Can specify in the config file a specific issue or range of issues to copy
- Auto-close source issues after copy
- Add a note with a link to the new issue created in the target project
- Use a custom link text template, like "Closed in favor or me/myotherproject#12"

## Usage

First, make sure you have valid GitLab account tokens for both source and destination GitLab installations. They are used
to access GitLab resources without authentication. GitLab private tokens are availble in "*Profile Settings* -> *Account*".

Now, write a `gitlab.yml` YAML config file to specify source and target projects, along with your GitLab account tokens:

```yaml
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

```yaml
from:
  url: https://gitlab.mydomain.com
  token: atoken
  project: namespace/project
  issues:
  - 15
  - 20-30
...
```

In order to copy all labels from one project to another (labels only, not issues), just append a `labelsOnly`
entry in the `from` section:

```yaml
from:
  url: https://gitlab.mydomain.com
  token: atoken
  project: namespace/project
  labelsOnly: true
to:
  url: https://gitlab.sameorotherdomain.com
  token: anothertoken
  project: namespace/otherproject
...
```

As of version `v0.6.6`, notes in issues can preserve original user ownership when copied. To do that, you need
to

- have tokens for all users involved
- add related users as members of the target project beforehand (with at least a *Reporter* permission)
- add a `users` entry into the `to` target section:

```yaml
...
to:
  url: https://gitlab.sameorotherdomain.com
  token: anothertoken
  project: namespace/otherproject
  users:
    bob: anothertoken
    alice: herowntoken
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

Ensure you have a working [Go](https://www.golang.org) 1.5+ installation then install [gb](http://getgb.io) to compile the project:
```
$ go get github.com/constabulary/gb/...
```

- To build the project, just run `make`. The program gets compiled into `bin/gitlab-copy`
- Cross-compile with `make buildall`
- Prepare distribution packages with `make dist`

## License

MIT. See `LICENSE` file.
