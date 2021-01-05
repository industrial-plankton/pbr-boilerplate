package Writing

func FormatErrors(errs []error) [][]interface{} {
	errors := make([][]interface{}, len(errs))
	for i := range errs {
		errors[i] = []interface{}{errs[i].Error()}
	}

	return errors
}
