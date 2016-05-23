# go-semver

[![Build Status](https://travis-ci.org/ceralena/go-semver.svg?branch=master)](https://travis-ci.org/ceralena/go-semver)

[![codecov](https://codecov.io/gh/ceralena/go-semver/branch/master/graph/badge.svg)](https://codecov.io/gh/ceralena/go-semver)

A [semver](http://semver.org) package for [go](https://golang.org).

See [godoc](https://godoc.org/github.com/ceralena/go-semver) for the API.

This package was designed to be as lightweight and simple as possible. It
provides a single `Version` type:

```go
type Version struct {
	Major      int
	Minor      int
	Patch      int
	Prerelease string
}
```

The type has methods for comparison with other Version values: `Equals`, `GreaterThan` and `LessThan`.

The package has a test suite with 100% coverage. To run the tests:

	go test -v -cover github.com/ceralena/go-semver

## Non-Standard Behaviour

The `Parse()` function supports incomplete version strings where the meaning can be
inferred. For example:

* `v2` -> `v2.0.0`
* `v2.4` -> `v2.4.0`

## Usage

See [godoc](https://godoc.org/github.com/ceralena/go-semver) for the API.

Assuming familiarity with [how to write Go
code](https://golang.org/doc/code.html):

To get the package, run `go get github.com/ceralena/go-semver`, or
just add it to the imports for your package:

```go
import "github.com/ceralena/go-semver"
```

To parse a version string:

```go
s := "v1.0.7-alpha"
v, err := semver.Parse(s)
if err != nil {
	panic(err)
}
fmt.Println(v)
```

To compare two versions as equal, greater or less than:

```go
a := semver.Version{1, 2, 3, "rc.1"}
b := semver.Version{4, 2, 3, "rc.1"}

a.Equals(b) // false
a.GreaterThan(b) // false
a.LessThen(b) // true
```

Printing the `Version` type will produce a standard semver string:

```go
a := semver.Version{4, 12, 1, "beta"}
fmt.Sprintf("%s", a) // v4.12.1-beta
```

Contributing
------------

All contributions and bug reports are welcome.

If you have any questions or issues, please use GitHub's issue tracker.

If you want to contribute, you can either send a pull request or [contact me on
twitter](https://twitter.com/ceralena).

License
-------

This package uses the MIT license. Generally speaking: you can do anything with this code. See the `LICENSE` file for the full text.
