package util

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	constanta "github.com/faithol1024/bgp-hackaton/lib/const"
	"github.com/google/uuid"
)

func GetEnv() string {
	environ := os.Getenv("TKPENV")
	if environ == "" {
		environ = constanta.EnvDevelopment
	}
	return environ
}

func GetAppName() string {
	return "campaign-engine"
}

func IsDevelopmentEnv() bool {
	return GetEnv() == constanta.EnvDevelopment
}

func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func StrintToInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

func GetLineOfCode(skip int) string {
	pc, _, line, _ := runtime.Caller(skip)
	details := fmt.Sprintf(
		"%s[%d]",
		runtime.FuncForPC(pc).Name(),
		line,
	)

	return details
}

func GetStringUUID() string {
	id := uuid.New()
	stringID := strings.Replace(id.String(), "-", "", -1)
	return stringID
}
