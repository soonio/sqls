package sqls

type Opt struct {
	ignore []string
	only   []string
	where  []string
	suffix []string
}

func (o *Opt) Ignore(column ...string) *Opt {
	o.ignore = append(o.ignore, column...)
	return o
}

func (o *Opt) Only(column ...string) *Opt {
	o.only = append(o.only, column...)
	return o
}
func (o *Opt) Where(condition ...string) *Opt {
	o.where = append(o.where, condition...)
	return o
}

func (o *Opt) Suffix(suffix ...string) *Opt {
	o.suffix = append(o.suffix, suffix...)
	return o
}
