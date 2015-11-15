package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/codegangsta/cli"
	"github.com/myshkin5/netspel/adapters/udp"
	"github.com/myshkin5/netspel/factory"
	"github.com/myshkin5/netspel/logs"
	"github.com/myshkin5/netspel/schemes/simple"
	"github.com/op/go-logging"
	"github.com/myshkin5/netspel/adapters/sse"
)

func init() {
	factory.WriterManager.RegisterType("netspel.adapters.udp.Writer", reflect.TypeOf(udp.Writer{}))
	factory.ReaderManager.RegisterType("netspel.adapters.udp.Reader", reflect.TypeOf(udp.Reader{}))

	factory.WriterManager.RegisterType("netspel.adapters.sse.Writer", reflect.TypeOf(sse.Writer{}))
	factory.ReaderManager.RegisterType("netspel.adapters.sse.Reader", reflect.TypeOf(sse.Reader{}))

	factory.SchemeManager.RegisterType("netspel.schemes.simple.Scheme", reflect.TypeOf(simple.Scheme{}))
}

func main() {
	app := cli.NewApp()
	app.Name = "netspel"
	app.Usage = "test network throughput with varying protocols"
	app.HideVersion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "configuration file",
		},
		cli.StringFlag{
			Name:  "scheme, s",
			Usage: "scheme type overriding the config file",
		},
		cli.StringFlag{
			Name:  "writer, w",
			Usage: "writer type overriding the config file",
		},
		cli.StringFlag{
			Name:  "reader, r",
			Usage: "reader type overriding the config file",
		},
		cli.StringSliceFlag{
			Name:  "config-string",
			Usage: "additional configuration <key>=<value> strings overriding the config file",
		},
		cli.StringSliceFlag{
			Name:  "config-int",
			Usage: "additional configuration <key>=<value> integers overriding the config file",
		},
		cli.StringFlag{
			Name:   "log-level, l",
			Usage:  "logging level",
			EnvVar: "INFO,DEBUG",
		},
	}
	app.Commands = []cli.Command{
		cli.Command{
			Name:    "write",
			Aliases: []string{"w"},
			Usage:   "write messages",
			Action: func(context *cli.Context) {
				write(context)
			},
		},
		cli.Command{
			Name:    "read",
			Aliases: []string{"r"},
			Usage:   "read messages",
			Action: func(context *cli.Context) {
				read(context)
			},
		},
	}

	app.RunAndExitOnError()
}

func write(context *cli.Context) {
	initLogs(context)

	config := config(context)
	scheme := scheme(config, context)

	writer, err := factory.CreateWriter(config.WriterType)
	if err != nil {
		panic(err)
	}

	err = writer.Init(config.Additional)
	if err != nil {
		cli.ShowAppHelp(context)
		panic(err)
	}

	scheme.RunWriter(writer)

	outputReport(scheme)
}

func read(context *cli.Context) {
	initLogs(context)

	config := config(context)
	scheme := scheme(config, context)

	reader, err := factory.CreateReader(config.ReaderType)
	if err != nil {
		cli.ShowAppHelp(context)
		panic(err)
	}

	err = reader.Init(config.Additional)
	if err != nil {
		panic(err)
	}

	scheme.RunReader(reader)

	outputReport(scheme)
}

func initLogs(context *cli.Context) {
	level, err := logging.LogLevel(context.GlobalString("log-level"))
	if err != nil {
		level = logging.INFO
	}

	logs.LogLevel.SetLevel(level, "netspel")
}

func config(context *cli.Context) factory.Config {
	configPath := context.GlobalString("config")
	var config factory.Config
	var err error
	if configPath == "" {
		config, err = factory.Parse([]byte("{}"))
	} else {
		config, err = factory.LoadFromFile(configPath)
	}
	if err != nil {
		cli.ShowAppHelp(context)
		panic(err)
	}

	schemeType := context.GlobalString("scheme")
	if schemeType != "" {
		config.SchemeType = schemeType
	}
	writerType := context.GlobalString("writer")
	if writerType != "" {
		config.WriterType = writerType
	}
	readerType := context.GlobalString("reader")
	if readerType != "" {
		config.ReaderType = readerType
	}

	for _, assignment := range context.GlobalStringSlice("config-string") {
		keyValue, err := parseAssignment(assignment)
		if err != nil {
			panic(err)
		}

		config.Additional.SetString(keyValue[0], keyValue[1])
	}
	for _, assignment := range context.GlobalStringSlice("config-int") {
		keyValue, err := parseAssignment(assignment)
		if err != nil {
			panic(err)
		}

		value, err := strconv.Atoi(keyValue[1])
		if err != nil {
			panic(err)
		}

		config.Additional.SetInt(keyValue[0], value)
	}

	return config
}

func parseAssignment(assignment string) ([]string, error) {
	keyValue := strings.Split(assignment, "=")
	if len(keyValue) != 2 {
		return []string{}, fmt.Errorf("Values must be of the form <key>=<value>, %s", assignment)
	}

	return keyValue, nil
}

func scheme(config factory.Config, context *cli.Context) factory.Scheme {
	scheme, err := factory.CreateScheme(config.SchemeType)
	if err != nil {
		cli.ShowAppHelp(context)
		panic(err)
	}

	err = scheme.Init(config.Additional)
	if err != nil {
		panic(err)
	}

	return scheme
}

func outputReport(scheme factory.Scheme) {
	bytesPerSec := ByteSize(scheme.ByteCount()) * ByteSize(time.Second) / ByteSize(scheme.RunTime().Nanoseconds())
	messagesPerSec := float64(scheme.MessagesPerRun()) * float64(time.Second) / float64(scheme.RunTime().Nanoseconds())

	fmt.Printf("\n\nByte count: %d\n", scheme.ByteCount())
	fmt.Printf("Rates: %s/s %.1f messages/s\n", bytesPerSec.String(), messagesPerSec)
	fmt.Printf("Error count: %d\n", scheme.ErrorCount())
	fmt.Printf("Run time: %s\n", scheme.RunTime().String())
	if scheme.FirstError() != nil {
		fmt.Printf("First error: %s\n", scheme.FirstError().Error())
	}
}
