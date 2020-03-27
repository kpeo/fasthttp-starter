package setup

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	stdlog "log"
	"reflect"
)

var (
	appName    = atomic.NewString("setup")
	appVersion = atomic.NewString("dev")
)

// Catch errors
func Catch(err error) {
	if err == nil {
		return
	}

	v := viper.New()

	log, logErr := NewLogger(NewLoggerConfig(v), &Settings{
		Name:    appName.Load(),
		Version: appVersion.Load(),
	})

	if logErr != nil {
		stdlog.Fatal(err)
	} else {
		log.Fatal("Can't run app", zap.Error(err))
	}
}

// CatchTrace catch errors with backtrace
// use it for debugging only
func CatchTrace(err error) {
	if err == nil {
		return
	}
	// dig into the source of the error
loop:
	for {
		var (
			ok bool
			v  = reflect.ValueOf(err)
			fn reflect.Value
		)

		switch {
		case v.Type().Kind() != reflect.Struct,
			!v.FieldByName("Reason").IsValid():
			break loop
		case v.FieldByName("Func").IsValid():
			fn = v.FieldByName("Func")
		}

		fmt.Printf("Func: %#v\nReason: %s\n\n", fn, err)

		if err, ok = v.FieldByName("Reason").Interface().(error); !ok {
			err = v.Interface().(error)
			break
		}
	}

	panic(err)
}
