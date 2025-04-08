package service

import "geolize/utilities/conf"

var (
	name string
	port int32
)

func init() {
	name, _ = conf.GetString("service", "name", "service")
	port, _ = conf.GetInt32("service", "port", 9000)
}

func GetPort() int32 {
	return port
}

func GetName() string {
	return name
}
