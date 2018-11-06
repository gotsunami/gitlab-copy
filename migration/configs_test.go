package migration

const cfg1 = `
from:
    url: https://gitlab.mydomain.com
    token: sourcetoken
    project: source/project
#    issues:
#    - 5
#    - 8-10
    labelsOnly: true
    # moveIssues: true
to:
    url: https://gitlab.mydomain.com
    token: desttoken
    project: dest/project
`

const cfg2 = `
from:
    url: https://gitlab.mydomain.com
    token: sourcetoken
    project: source/project
to:
    url: https://gitlab.mydomain.com
    token: desttoken
    project: dest/project
`

const cfg3 = `
from:
    url: https://gitlab.mydomain.com
    token: sourcetoken
    project: source/project
    milestonesOnly: true
to:
    url: https://gitlab.mydomain.com
    token: desttoken
    project: dest/project
`
