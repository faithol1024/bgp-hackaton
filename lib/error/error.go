package ers

import (
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
