package filter

type Filter struct {
	Conditions []Condition
}

type Condition struct {
	Key   string
	Value interface{}
}

func MakeFilter(conditions ...Condition) Filter {
	return Filter{
		Conditions: conditions,
	}
}

func (f *Filter) AddCondition(condition Condition) {
	f.Conditions = append(f.Conditions, condition)
}
