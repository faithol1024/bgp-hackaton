package ers

import (
	"errors"
	"strings"

	"github.com/faithol1024/bgp-hackaton/lib/util"
)

type ErrorWithDetails struct {
	Err            error
	TrueErrMessage string
	Traces         []string
}

func (h ErrorWithDetails) Error() string {
	return h.Err.Error()
}

// Add line of code for easy error trace
func ErrorAddTrace(err interface{}) error {
	if _, ok := err.(error); !ok {
		return nil
	}

	// else, create new error with details
	errd := ErrorWithDetails{
		Err:    err.(error),
		Traces: []string{getCleanTrace(util.GetLineOfCode(2))},
	}

	return errd
}

func getCleanTrace(traceinput string) (trace string) {
	traces := strings.Split(traceinput, "github.com/faithol1024/bgp-hackaton")
	trace = traces[0]
	if len(traces) > 1 {
		trace = traces[1]
	}

	return trace
}

func ErrorGetTrace(err interface{}) []string {
	if _, ok := err.(error); !ok {
		return []string{}
	}

	if errd, ok := err.(ErrorWithDetails); ok {
		return errd.Traces
	}

	return []string{}
}

func IsMatchError(err1 error, err2 error) bool {
	if err1 == nil && err2 == nil {
		return true
	}

	if err1 == nil {
		err1 = errors.New("nil")
	}

	if err2 == nil {
		err2 = errors.New("nil")
	}

	// for now comparing the message only, because if comparing errors will panic
	if err1.Error() == err2.Error() {
		return true
	}

	return false
}
