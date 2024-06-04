package storage

import (
	"path/filepath"
	"strings"
)

// Prefix contains the attributes which identify the pathway to S3 objects. It does not include the bucket
type Prefix struct {
	Base     string
	Bible    string
	Fileset  string
	Filename string
}

func NewPrefix(base string, bible string, fileset string, filename string) Prefix {
	return Prefix{Base: base, Bible: bible, Fileset: fileset, Filename: filename}
}

func Parse(prefix string) Prefix {
	var base, bible, fileset, filename string
	parts := strings.Split(prefix, "/")
	ct := len(parts)
	if ct > 0 {
		base = parts[0]
	}
	if ct > 1 {
		bible = parts[1]
	}
	if ct > 2 {
		fileset = parts[2]
	}
	if ct > 3 {
		filename = parts[3]
	}
	return Prefix{Base: base, Bible: bible, Fileset: fileset, Filename: filename}
}

// String returns a string representation of Prefix
func (p Prefix) String() string {
	return filepath.Join(p.Base, p.Bible, p.Fileset, p.Filename)
}

// Compare performs piece comparison of each of the prefix attributes
func Compare(a Prefix, b Prefix) int {
	if a.Base == b.Base {
		if a.Bible == b.Bible {
			if a.Fileset == b.Fileset {
				if a.Filename == b.Filename {
					return 0
				}
				if a.Filename < b.Filename {
					return -1
				}
				return +1
			}
			if a.Fileset < b.Fileset {
				return -1
			}
			return +1

		}
		if a.Bible < b.Bible {
			return -1
		}
		return +1
	}
	if a.Base < b.Base {
		return -1
	}
	return +1
}
