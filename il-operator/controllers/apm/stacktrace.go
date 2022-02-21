package apm

import "github.com/pkg/errors"

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func stackTrace(err error) []uintptr {
	st := deepestStackTrace(err)
	if st == nil {
		return nil
	}
	return transformStackTrace(st)
}

func deepestStackTrace(err error) errors.StackTrace {
	var last stackTracer
	for err != nil {
		if err, ok := err.(stackTracer); ok {
			last = err
		}
		cause, ok := err.(interface {
			Cause() error
		})
		if !ok {
			break
		}
		err = cause.Cause()
	}

	if last == nil {
		return nil
	}
	return last.StackTrace()
}

func transformStackTrace(orig errors.StackTrace) []uintptr {
	st := make([]uintptr, len(orig))
	for i, frame := range orig {
		st[i] = uintptr(frame)
	}
	return st
}
