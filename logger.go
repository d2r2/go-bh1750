package bh1750

import logger "github.com/d2r2/go-logger"

// You can manage verbosity of log output
// in the package by changing last parameter value.
var lg = logger.NewPackageLogger("bh1750",
	logger.DebugLevel,
	// logger.InfoLevel,
)
