package nosqldb

type Filter int

const (
	PersonDoc Filter = iota
	DefaultsDoc

	PersonType   string = "person"
	DefaultsType string = "defaults"
)

func (t Filter) Tuple() (string, string, interface{}) {
	switch t {
	case PersonDoc:
		return "type", "==", PersonType
	case DefaultsDoc:
		return "type", "==", DefaultsType
	default:
		return "", "", nil
	}
}
