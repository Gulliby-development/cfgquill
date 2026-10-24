package config

// Entry is one configuration variable and its parsed value, if present.
type Entry struct {
	Line     int
	Name     string
	Var      string
	Key      string
	Value    string
	HasValue bool
}

// Header is one section header and the number of completed entries after it.
type Header struct {
	Line  int
	Name  string
	Count int
}

// File is the ordered configuration model for a parsed text.
type File struct {
	Entries []Entry
	Headers []Header
}

// Find returns entries whose stored key exactly matches canon, in file order.
func (f File) Find(canon string) []Entry {
	var found []Entry
	for _, entry := range f.Entries {
		if entry.Key == canon {
			found = append(found, entry)
		}
	}
	return found
}
