# releasedver

[![Build Status](https://travis-ci.org/ksoichiro/releasedver.svg?branch=master)](https://travis-ci.org/ksoichiro/releasedver)

`releasedver` forces to use go modules with released version in go.mod file.

- OK: `example.com/foo v0.1.0`
- NG: `example.com/foo v0.0.0-20200051612345-fbaf1a3ff141`

It might be useful in these use cases:

- You don't want to allow to use the commit without tag because it is unstable due to their development workflow.
- You don't want to allow to use the commit without tag because it might be deleted in the future if it's just a temporary feature branch.

## How to use

```
go get -u github.com/ksoichiro/releasedver/cmd/releasedver
go vet -vettool=`which releasedver` -paths=golang.org/x/tools,github.com/gostaticanalysis/modfile ./...
```

output:

```
# github.com/ksoichiro/releasedver
./go.mod:6:1: github.com/gostaticanalysis/modfile must use released version
./go.mod:8:1: golang.org/x/tools must use released version
```

## Limitations

## Go 1.12 support

For Go 1.12, you must set `GO111MODULE=on` to make the linter to work.

## go vet

When you execute this linter by `go vet` and if your project has a structure
that the project root directory has go.mod file but does not have any Go files,
then you should use `-root` flag to specify the root directory path.
Otherwise the linter cannot find go.mod file.

```
go vet -vettool=`which releasedver` -paths=golang.org/x/tools,github.com/gostaticanalysis/modfile -root=. ./...
```

Currently, the linter reports duplicate issues when executed by `go vet`,
so if your project has this structure it's better to use the binary directly.

```
releasedver -paths=golang.org/x/tools,github.com/gostaticanalysis/modfile ./...
```
