// package semver provides a type representing a semantic version, and
// facilities for parsing, serialisation and comparison.
//
// See http://semver.org for more information on semantic versioning.
//
// This package expands on the specification: a partial version string like
// "v2" or "v2.0" is considered valid, and expanded to "v2.0.0".
//
//  To parse a version string:
//
//    s := "v1.0.7-alpha"
//    v, err := semver.Parse(s)
//    if err != nil {
//  	panic(err)
//    }
//    fmt.Println(v)
//
// Visit godoc.org/github.com/DanielBryan/semver for the full package API.
package semver

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strconv"
)

// A Version is a parsed semver version string.
//
// If only a partial version was specified, the missing parts will be -1.
type Version struct {
	Major      int
	Minor      int
	Patch      int
	Prerelease string
}

// Returns the standard string representation of the Version value.
//
// For example, given this Version value:
//
//  Version{Major: 1, Minor: 2, Patch: 3, Prerelease: "beta1"}
//
// The following string is produced:
//
//  v1.2.3-beta1
func (v Version) String() string {
	if len(v.Prerelease) > 0 {
		return fmt.Sprintf("v%d.%d.%d-%s", v.Major, v.Minor, v.Patch, v.Prerelease)
	} else {
		return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
	}
}

// Returns true if v is a higher version than o.
func (v Version) GreaterThan(o Version) bool {
	if v.Major > o.Major {
		return true
	} else if v.Minor > o.Minor {
		return true
	} else if v.Patch > o.Patch {
		return true
	} else if v.Prerelease == o.Prerelease {
		return false
	}

	sl := []string{v.Prerelease, o.Prerelease}
	sort.Strings(sl)
	if sl[0] == v.Prerelease {
		return false
	}
	return true
}

// Returns true if v is a lesser version than o.
func (v Version) LessThan(o Version) bool {
	return !v.Equals(o) && !v.GreaterThan(o)
}

// Returns true if v and o are the same version.
func (v Version) Equals(o Version) bool {
	return v.Major == o.Major && v.Minor == o.Minor && v.Patch == o.Patch && v.Prerelease == o.Prerelease
}

// States for parsing state machine.
// The states are ordered - parsing can only advance forwards.
// We can only advance forwards.
type parseState int

const (
	atStart parseState = iota
	foundV
	foundMajor
	foundMinor
	foundPatch
	foundPrerelease
)

var (
	EmptyVersion   = errors.New("Empty version string")
	IllegalVersion = errors.New("Illegal version string")
)

// Parse a version value from its string representation.
//
// The error value will either be nil, EmptyVersion or IllegalVersion.
func Parse(s string) (Version, error) {
	var (
		v     Version                          // Version object we'll return
		state parseState             = atStart // current parsing state
		pos   int                              // pointer into the string
		buf   = bytes.NewBuffer(nil)           // container for temporary state while we loop
		err   error
	)

	if len(s) == 0 {
		return v, EmptyVersion
	}

	// Loop until we find an error or we've finished parsing the string

	for pos < len(s) {
		switch state {
		case atStart:
			if s[pos] == 'v' {
				pos = pos + 1
				state = foundV
			} else {
				return v, IllegalVersion
			}
		case foundV:
			var maj int
			if maj, pos, err = readNextNum(s, pos, buf); err != nil {
				return v, err
			}
			v.Major = maj
			state = foundMajor
		case foundMajor:
			var minor int
			if minor, pos, err = readNextNum(s, pos, buf); err != nil {
				return v, err
			}
			v.Minor = minor
			state = foundMinor
		case foundMinor:
			var patch int
			if patch, pos, err = readNextNum(s, pos, buf); err != nil {
				return v, err
			}
			v.Patch = patch
			state = foundPatch
		case foundPatch:
			v.Prerelease = s[pos:]
			pos = len(s)
			state = foundPrerelease
		}
	}

	if state < foundMajor {
		// At minimum we need a major version
		return v, IllegalVersion
	}

	return v, nil
}

// Read the next version number from this cursor in the string.
// buf should be an empty bytes.Buffer. The buffer will be automatically reset.
//
// Reads until a period, hyphen or the end of the string.
//
// Returns the version number, the new cursor point and any error.
func readNextNum(s string, curs int, buf *bytes.Buffer) (int, int, error) {
	defer buf.Reset()
	for ; curs < len(s) && s[curs] != '.' && s[curs] != '-'; curs += 1 {
		buf.WriteByte(s[curs])
	}
	i, err := strconv.Atoi(buf.String())
	if err != nil {
		return -1, curs, IllegalVersion
	}
	return i, curs + 1, nil
}
