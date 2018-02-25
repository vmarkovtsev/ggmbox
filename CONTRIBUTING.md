# Contributing Guidelines

ggmbox project is [MIT licensed](LICENSE) and accepts
contributions via GitHub pull requests.  This document outlines some of the
conventions on development workflow, commit message formatting, contact points,
and other resources to make it easier to get your contribution accepted.

## How to Contribute

Pull Requests (PRs) are the main and exclusive way to contribute to the official ggmbox project.
In order for a PR to be accepted it needs to pass a list of requirements:

- All the tests pass.
- Python code is formatted according to [![PEP8](https://img.shields.io/badge/code%20style-pep8-orange.svg)](https://www.python.org/dev/peps/pep-0008/).
- Go code is idiomatic and passes `gofmt` and `golint` without any warnings.
- If the PR is a bug fix, it has to include a new unit test that fails before the patch is merged.
- If the PR is a new feature, it has to come with a suite of unit tests, that tests the new functionality.
- In any case, all the PRs have to pass the personal evaluation of at least one of the [maintainers](MAINTAINERS.md).


### Format of the commit message

The commit summary must start with a capital letter and with a verb in present tense. No dot in the end.

```
Add a feature
Remove unused code
Fix a bug
```

Every commit details should describe what was changed, under which context and, if applicable, the GitHub issue it relates to.
