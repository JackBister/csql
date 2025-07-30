package csql

type Options struct {
	PrintOps   bool
	PrintTypes bool
	Separator  string
	Skip       int
}

func NewOptions() Options {
	return Options{
		PrintOps:   false,
		PrintTypes: false,
		Separator:  ",",
		Skip:       0,
	}
}
