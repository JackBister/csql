package csql

type Options struct {
	PrintTypes bool
	Separator  string
	Skip       int
}

func NewOptions() Options {
	return Options{
		PrintTypes: false,
		Separator:  ",",
		Skip:       0,
	}
}
