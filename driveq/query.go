// Package driveq
package driveq

import (
	"strings"
	"time"
)

type Query interface {
	String() string
}

// --- grouping with curly braces ---

type group struct {
	q Query
}

func (g group) String() string {
	return "{ " + g.q.String() + " }"
}

func Group(q Query) Query {
	return group{q}
}

// --- raw escape hatch ---

type raw string

func (r raw) String() string {
	return string(r)
}

// --- quoting rules ---

func quote(v string) string {
	if v == "true" || v == "false" {
		return v
	}
	if strings.HasPrefix(v, "'") && strings.HasSuffix(v, "'") {
		return v
	}
	return "'" + strings.ReplaceAll(v, "'", "\\'") + "'"
}

// --- time helpers ---

func RFC3339(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}
