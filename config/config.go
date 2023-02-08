package config

import "os"

var LogFile = os.Getenv("LOG_FILE")
