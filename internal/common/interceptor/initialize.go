package interceptor

import "buf.build/go/protovalidate"

var (
	Validator protovalidate.Validator
)

func Initialize() error {
	var err error
	Validator, err = protovalidate.New()
	if err != nil {
		return err
	}
	return nil
}
