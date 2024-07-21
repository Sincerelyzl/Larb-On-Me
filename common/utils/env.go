package utils

import (
	"strconv"
	"syscall"
	"time"
)

func EnvString(key string, defaultValue string) string {
	value, exist := syscall.Getenv(key)
	if !exist {
		return defaultValue
	}
	return value
}

func EnvInt(key string, defaultValue int) int {
	value, exist := syscall.Getenv(key)
	if !exist {
		return defaultValue
	}

	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return valueInt
}

func EnvBool(key string, defaultValue bool) bool {
	value, exist := syscall.Getenv(key)
	if !exist {
		return defaultValue
	}

	valueBool, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return valueBool
}

func EnvDuration(key string, defaultValue string) time.Duration {
	// example: 5s, 10m, 1h
	value, exist := syscall.Getenv(key)
	if !exist {
		duration, err := time.ParseDuration(defaultValue)
		if err != nil {
			return 0
		}
		return duration
	}

	valueDuration, err := time.ParseDuration(value)
	if err != nil {
		duration, err := time.ParseDuration(defaultValue)
		if err != nil {
			return 0
		}
		return duration
	}

	return valueDuration
}
