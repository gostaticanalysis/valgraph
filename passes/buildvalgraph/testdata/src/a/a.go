package a

import (
	"errors"
	"fmt"
)

type pkgErr struct {
	err error
}

func (e *pkgErr) Error() string {
	return fmt.Sprintf("pkg error: %s", e.err.Error())
}

func newError(n int) error {
	if n%2 == 0 {
		return nil
	}
	return errors.New("error")
}

func noError() {
	println("no error")
}

func directolyError(n int) error {
	err := newError(n)
	if err != nil {
		return err
	}
	return nil
}

func wrapError(err error) error {
	return &pkgErr{
		err: err,
	}
}

func wrappingError(n int) error {
	err := newError(n)
	if err != nil {
		return wrapError(err)
	}
	return nil
}

func wrappingError2(n int) error {
	err := newError(n)
	if err != nil {
		return fmt.Errorf("wrappingError2: %w", err)
	}
	return nil
}
