package driveq

import "strings"

type logical struct {
	op    string
	terms []Query
}

func (l logical) String() string {
	parts := make([]string, 0, len(l.terms))
	for _, t := range l.terms {
		parts = append(parts, t.String())
	}
	return strings.Join(parts, " "+l.op+" ")
}

func And(q ...Query) Query {
	return logical{"and", q}
}

func Or(q ...Query) Query {
	return logical{"or", q}
}
