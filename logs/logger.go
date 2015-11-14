package logs

import (
	"os"

	gologging "github.com/op/go-logging"
)

var (
	Logger   *gologging.Logger
	LogLevel gologging.LeveledBackend
)

func init() {
	Logger = gologging.MustGetLogger("netspel")
	format := gologging.MustStringFormatter("%{time:2006-01-02T15:04:05.000000Z} %{level} %{message}")
	backend := gologging.NewLogBackend(os.Stdout, "", 0)
	backendFormatter := gologging.NewBackendFormatter(backend, format)
	LogLevel = gologging.AddModuleLevel(backendFormatter)
	gologging.SetBackend(LogLevel)
}
