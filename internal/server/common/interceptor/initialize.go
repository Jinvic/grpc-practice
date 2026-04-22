package interceptor

func Initialize() error {
	if err := InitValidateInterceptor(); err != nil {
		return err
	}
	if err := InitLoggingInterceptor(); err != nil {
		return err
	}
	return nil
}
