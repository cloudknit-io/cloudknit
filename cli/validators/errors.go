package validators

type InvalidEnvironmentComponent struct {
	Component string
	Err error
}

func (e InvalidEnvironmentComponent) Error() string {
	return e.Err.Error()
}
