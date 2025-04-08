package maxmind

import (
	"geolize/utilities/conf"
)

const (
	versionFilePath = "data/version"
	dbFolder        = "data/db/"
	dbHistories     = "data/histories/"
)

var (
	db, _ = conf.GetString("geolize", "db", "GeoLite2-City.mmdb")
)
