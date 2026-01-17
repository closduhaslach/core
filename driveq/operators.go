package driveq

type Condition struct {
	Field Field
	Op    string
	Value string
}

func (c Condition) String() string {
	return string(c.Field) + " " + c.Op + " " + quote(c.Value)
}

func Eq(f Field, v string) Query {
	return Condition{f, "=", v}
}

func Contains(f Field, v string) Query {
	return Condition{f, "contains", v}
}

func Gt(f Field, v string) Query {
	return Condition{f, ">", v}
}

func Gte(f Field, v string) Query {
	return Condition{f, ">=", v}
}

func Lt(f Field, v string) Query {
	return Condition{f, "<", v}
}

func Lte(f Field, v string) Query {
	return Condition{f, "<=", v}
}

func In(f Field, v string) Query {
	return raw("'" + v + "' in " + string(f))
}

func Has(field Field, email string) Query {
	return raw(string(field) + " has '" + email + "'")
}
