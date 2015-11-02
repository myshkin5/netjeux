package main

import (
	"reflect"

	"netspel/adapters/udp"
	"netspel/factory"
	"netspel/schemes/simple"

	"github.com/codegangsta/cli"
)

func init() {
	factory.WriterManager.RegisterType("netspel.adapters.udp.Writer", reflect.TypeOf(udp.Writer{}))
	factory.ReaderManager.RegisterType("netspel.adapters.udp.Reader", reflect.TypeOf(udp.Reader{}))
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
		cli.StringSliceFlag{
			Name:  "config-string, s",
			Usage: "additional configuration <key>=<value> strings overriding the config file",
		},
		cli.StringSliceFlag{
			Name:  "config-int, i",
			Usage: "additional configuration <key>=<value> integers overriding the config file",
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
	config := config(context)
	scheme := scheme(config, context)

	writer, err := factory.CreateWriter(config.WriterType)
	if err != nil {
		panic(err)
	}

	err = writer.Init(config)
	if err != nil {
		panic(err)
	}

	scheme.RunWriter(writer)
}

func read(context *cli.Context) {
	config := config(context)
	scheme := scheme(config, context)

	reader, err := factory.CreateReader(config.ReaderType)
	if err != nil {
		panic(err)
	}

	err = reader.Init(config)
	if err != nil {
		panic(err)
	}

	scheme.RunReader(reader)
}

func config(context *cli.Context) factory.Config {
	config, err := factory.LoadFromFile(context.GlobalString("config"))
	if err != nil {
		cli.ShowAppHelp(context)
		panic(err)
	}

	for _, assignment := range context.GlobalStringSlice("config-string") {
		config.ParseAndSetAdditionalString(assignment)
	}
	for _, assignment := range context.GlobalStringSlice("config-int") {
		config.ParseAndSetAdditionalInt(assignment)
	}

	return config
}

func scheme(config factory.Config, context *cli.Context) factory.Scheme {
	scheme, err := factory.CreateScheme(config.SchemeType)
	if err != nil {
		panic(err)
	}

	err = scheme.Init(config)
	if err != nil {
		panic(err)
	}

	return scheme
}
