package common

import (
	"github.com/modelbus/one-api-pro/common/env"
)

var UsingSQLite = false
var UsingPostgreSQL = false
var UsingMySQL = false

var SQLitePath = "one-api-pro.db"
var SQLiteBusyTimeout = env.Int("SQLITE_BUSY_TIMEOUT", 3000)
