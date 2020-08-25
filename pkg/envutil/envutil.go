package envutil

import (
	"fmt"
	"os"
	"strconv"
)

func MustGetEnvStr(env string) string {
	if val := GetEnvStr(env); len(val) > 0 {
		return val
	}
	panic(fmt.Errorf("env %s is not set", env))
}

func GetEnvStr(env string) string {
	return os.Getenv(env)
}

func GetEnvStrOrDefault(env string, defaultVal string) string {
	valStr := os.Getenv(env)
	if len(valStr) == 0 {
		return defaultVal
	}
	return valStr
}

func GetEnvUintOrDefault(env string, defaultVal uint64) uint64 {
	valStr := os.Getenv(env)
	if val, err := strconv.ParseUint(valStr, 10, 64); err == nil {
		return val
	}
	return defaultVal
}
