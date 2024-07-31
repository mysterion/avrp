package logg

import (
	"log"

	"github.com/mysterion/avrp/internal/utils"
)

func Info(s ...any) {
	var v []any
	v = append(v, "INFO: ")
	v = append(v, s...)
	log.Println(v...)
}

func Warn(s ...any) {
	var v []any
	v = append(v, "WARNING: ")
	v = append(v, s...)
	log.Println(v...)
}

func Error(s ...any) {
	var v []any
	v = append(v, "ERROR: ")
	v = append(v, s...)
	log.Println(v...)
}

func Debug(s ...any) {
	if utils.DEV {
		var v []any
		v = append(v, "DEBUG: ")
		v = append(v, s...)
		log.Println(v...)
	}
}
