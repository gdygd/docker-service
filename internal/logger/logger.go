package logger

import (
	"github.com/gdygd/goglib"
)

// //---------------------------------------------------------------------------
// // Log
// //---------------------------------------------------------------------------
var Log *goglib.OLog2 = goglib.InitLogEnv("./log", "delivery", 1) // level 1~9, auth-service
