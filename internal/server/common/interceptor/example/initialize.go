package exampleinterceptor

func Initialize() error {
	if err := InitValidateInterceptor(); err != nil {
		return err
	}
	return nil
}
