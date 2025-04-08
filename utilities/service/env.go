package service

import "os"

func GetEnv() string {
	env := os.Getenv("env")
	if len(env) < 1 {
		env = "dev"
	}

	return env
}

func IsDev() bool {
	return GetEnv() == "dev"
}

func IsProd() bool {
	return GetEnv() == "prod"
}

func IsTest() bool {
	return GetEnv() == "test"
}
