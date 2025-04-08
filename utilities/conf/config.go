package conf

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
)

var (
	env string
	f   *ini.File
)

func init() {
	env = os.Getenv("env")
	if len(env) < 1 {
		env = "dev"
	}

	var err error
	f, err = ini.Load(fmt.Sprintf("%s.ini", env))
	if err != nil {
		panic(fmt.Sprintf("missing %s.ini to start", env))
	}
}

func GetString(group string, key string, fallback string) (string, error) {
	if f != nil {
		k, err := f.Section(group).GetKey(key)
		if err != nil {
			return fallback, nil
		}

		return k.String(), nil
	}

	return "", fmt.Errorf("missing %s.ini to start", env)
}

func GetInt32(group string, key string, fallback int32) (int32, error) {
	if f != nil {
		k, err := f.Section(group).GetKey(key)
		if err != nil {
			return fallback, nil
		}

		v, err := k.Int64()
		if err != nil {
			return fallback, nil
		}

		return int32(v), nil
	}

	return 0, fmt.Errorf("missing %s.ini to start", env)
}

func GetBool(group string, key string, fallback bool) (bool, error) {
	if f != nil {
		k, err := f.Section(group).GetKey(key)
		if err != nil {
			return fallback, nil
		}

		v, err := k.Bool()
		if err != nil {
			return fallback, nil
		}

		return v, nil
	}

	return false, fmt.Errorf("missing %s.ini to start", env)
}

func GetStringSlice(group string, key string) ([]string, error) {
	if f != nil {
		k, err := f.Section(group).GetKey(key)
		if err != nil {
			return nil, nil
		}

		return k.Strings(","), nil
	}

	return nil, nil
}
