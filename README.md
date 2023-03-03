
# GitLab Copy

[![Build Status](https://travis-ci.org/gotsunami/gitlab-copy.svg?branch=master)](https://travis-ci.org/gotsunami/gitlab-copy)

`gitlab-copy` is a simple tool for copying issues/labels/milestones/notes from one GitLab project to another, possibly running on different GitLab instances.

By default, `gitlab-copy` won't copy anything until told **explicitly** to do so on the command line: running it will do nothing but showing some statistics.

**Note**: GitLab 8.6 introduced the ability [to move an issue to another project](https://about.gitlab.com/2016/03/22/gitlab-8-6-released/), but on the same GitLab installation only. `gitlab-copy` can still prove valuable to move issues between projects on different GitLab hosts and to perform batch operations from the command line (see the feature list below).

## Download

Installing `gitlab-copy` is very easy since it comes as a static binary with no dependencies. Just [grab a compiled version](https://github.com/gotsunami/gitlab-copy/releases/latest) for your system (or have a look at the **Compile From Source** section).

## Features

The following features are available:

- Support for GitLab instances with self-signed TLS certificates by using the `-k` CLI flag (since `v0.8`)
- Support for different GitLab hosts/instances (since `v0.8`)
- Copy milestones if not existing on target (use `milestonesOnly` to copy milestones only, see below)
- Copy all source labels on target (use `labelsOnly` to copy labels only, see below)
- Copy issues if not existing on target (by title)
- Apply closed status on issues, if any
- Set issue's assignee (if user exists) and milestone, if any
- Copy notes (attached to issues), preserving user ownership
- Can specify in the config file a specific issue or range of issues to copy
- Auto-close source issues after copy
- Add a note with a link to the new issue created in the target project
- Use a custom link text template, like "Closed in favor or me/myotherproject#12"

## Getting Started

Here are some instructions to get started. First make sure you have valid GitLab account tokens for both source and destination GitLab installations. They are used to access GitLab resources without authentication. GitLab private tokens are availble in "*Profile Settings* -> *Account*".

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

That's it. You may want to run the program now. See the section below.

## Run it!

Now grab some project stats by running
```
$ ./gitlab-copy gitlab.yml
```

If everything looks good, run the same command, this time with the `-y` flag to effectively copie issues between GitLab
instances (they can be the same):
```
$ ./gitlab-copy -y gitlab.yml
```

If one of the GitLab instances uses a self-signed TLS certificate, use the `-k` flag (available in `v0.8`) to skip the TLS verification process:

```
$ ./gitlab-copy -k -y gitlab.yml
```

## More Features

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

In order to copy all milestones only, just add a `milestonesOnly` entry in the `from` section:
```yaml
from:
  url: https://gitlab.mydomain.com
  token: atoken
  project: namespace/project
  milestonesOnly: true
...
```

Notes in issues can preserve original user ownership when copied. To do that, you need
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

## Compile From Source

Ensure you have a working [Go](https://www.golang.org) 1.18+ installation then:
```
$ go install github.com/gotsunami/gitlab-copy/cmd/gitlab-copy@latest
```

- The program gets compiled into `bin/gitlab-copy`
- Cross-compile with `make buildall`
- Prepare distribution packages with `make dist`

## Donate

If you like this tool and want to support its development, a donation would be greatly appreciated!

It's not about the amount at all: making a donation boosts the motivation to work on a project. Thank you very much if you can give anything.

Monero address: `88uoutKJS2w3FfkKyJFsNwKPHzaHfTAo6LyTmHSAoQHgCkCeR8FUG4hZ8oD4fnt8iP7i1Ty72V6CLMHi1yUzLCZKHU1pB7c`

![My monero address](qr-donate.png)

## License

MIT. See `LICENSE` file.
