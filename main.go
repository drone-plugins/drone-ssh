package main

import (
	"log"
	"os"
	"time"

	"github.com/appleboy/easyssh-proxy"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli/v2"
)

// Version set at compile-time
var Version string

func main() {
	// Load env-file if it exists first
	if filename, found := os.LookupEnv("PLUGIN_ENV_FILE"); found {
		_ = godotenv.Load(filename)
	}

	app := cli.NewApp()
	app.Name = "Drone SSH"
	app.Usage = "Executing remote ssh commands"
	app.Copyright = "Copyright (c) 2019 Bo-Yi Wu"
	app.Authors = []*cli.Author{
		{
			Name:  "Bo-Yi Wu",
			Email: "appleboy.tw@gmail.com",
		},
	}
	app.Action = run
	app.Version = Version
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "ssh-key",
			Usage:   "private ssh key",
			EnvVars: []string{"PLUGIN_SSH_KEY", "PLUGIN_KEY", "SSH_KEY", "KEY", "INPUT_KEY"},
		},
		&cli.StringFlag{
			Name:    "ssh-passphrase",
			Usage:   "The purpose of the passphrase is usually to encrypt the private key.",
			EnvVars: []string{"PLUGIN_SSH_PASSPHRASE", "PLUGIN_PASSPHRASE", "SSH_PASSPHRASE", "PASSPHRASE", "INPUT_PASSPHRASE"},
		},
		&cli.StringFlag{
			Name:    "key-path,i",
			Usage:   "ssh private key path",
			EnvVars: []string{"PLUGIN_KEY_PATH", "SSH_KEY_PATH", "INPUT_KEY_PATH"},
		},
		&cli.StringFlag{
			Name:    "username,user,u",
			Usage:   "connect as user",
			EnvVars: []string{"PLUGIN_USERNAME", "PLUGIN_USER", "SSH_USERNAME", "USERNAME", "INPUT_USERNAME"},
			Value:   "root",
		},
		&cli.StringFlag{
			Name:    "password,P",
			Usage:   "user password",
			EnvVars: []string{"PLUGIN_PASSWORD", "SSH_PASSWORD", "PASSWORD", "INPUT_PASSWORD"},
		},
		&cli.StringSliceFlag{
			Name:    "host,H",
			Usage:   "connect to host",
			EnvVars: []string{"PLUGIN_HOST", "SSH_HOST", "HOST", "INPUT_HOST"},
		},
		&cli.IntFlag{
			Name:    "port,p",
			Usage:   "connect to port",
			EnvVars: []string{"PLUGIN_PORT", "SSH_PORT", "PORT", "INPUT_PORT"},
			Value:   22,
		},
		&cli.BoolFlag{
			Name:    "sync",
			Usage:   "sync mode",
			EnvVars: []string{"PLUGIN_SYNC", "SYNC", "INPUT_SYNC"},
		},
		&cli.DurationFlag{
			Name:    "timeout,t",
			Usage:   "connection timeout",
			EnvVars: []string{"PLUGIN_TIMEOUT", "SSH_TIMEOUT", "TIMEOUT", "INPUT_TIMEOUT"},
			Value:   30 * time.Second,
		},
		&cli.DurationFlag{
			Name:    "command.timeout,T",
			Usage:   "command timeout",
			EnvVars: []string{"PLUGIN_COMMAND_TIMEOUT", "SSH_COMMAND_TIMEOUT", "COMMAND_TIMEOUT", "INPUT_COMMAND_TIMEOUT"},
			Value:   10 * time.Minute,
		},
		&cli.StringSliceFlag{
			Name:    "script,s",
			Usage:   "execute commands",
			EnvVars: []string{"PLUGIN_SCRIPT", "SSH_SCRIPT", "SCRIPT"},
		},
		&cli.StringFlag{
			Name:    "script.string",
			Usage:   "execute single commands for github action",
			EnvVars: []string{"INPUT_SCRIPT"},
		},
		&cli.BoolFlag{
			Name:    "script.stop",
			Usage:   "stop script after first failure",
			EnvVars: []string{"PLUGIN_SCRIPT_STOP", "STOP", "INPUT_SCRIPT_STOP"},
		},
		&cli.StringFlag{
			Name:    "proxy.ssh-key",
			Usage:   "private ssh key of proxy",
			EnvVars: []string{"PLUGIN_PROXY_SSH_KEY", "PLUGIN_PROXY_KEY", "PROXY_SSH_KEY", "INPUT_PROXY_KEY"},
		},
		&cli.StringFlag{
			Name:    "proxy.ssh-passphrase",
			Usage:   "The purpose of the passphrase is usually to encrypt the private key.",
			EnvVars: []string{"PLUGIN_PROXY_SSH_PASSPHRASE", "PLUGIN_PROXY_PASSPHRASE", "PROXY_SSH_PASSPHRASE", "PROXY_PASSPHRASE", "INPUT_PROXY_PASSPHRASE"},
		},
		&cli.StringFlag{
			Name:    "proxy.key-path",
			Usage:   "ssh private key path of proxy",
			EnvVars: []string{"PLUGIN_PROXY_KEY_PATH", "PROXY_SSH_KEY_PATH", "INPUT_PROXY_KEY_PATH"},
		},
		&cli.StringFlag{
			Name:    "proxy.username",
			Usage:   "connect as user of proxy",
			EnvVars: []string{"PLUGIN_PROXY_USERNAME", "PLUGIN_PROXY_USER", "PROXY_SSH_USERNAME", "INPUT_PROXY_USERNAME"},
			Value:   "root",
		},
		&cli.StringFlag{
			Name:    "proxy.password",
			Usage:   "user password of proxy",
			EnvVars: []string{"PLUGIN_PROXY_PASSWORD", "PROXY_SSH_PASSWORD", "INPUT_PROXY_PASSWORD"},
		},
		&cli.StringFlag{
			Name:    "proxy.host",
			Usage:   "connect to host of proxy",
			EnvVars: []string{"PLUGIN_PROXY_HOST", "PROXY_SSH_HOST", "INPUT_PROXY_HOST"},
		},
		&cli.StringFlag{
			Name:    "proxy.port",
			Usage:   "connect to port of proxy",
			EnvVars: []string{"PLUGIN_PROXY_PORT", "PROXY_SSH_PORT", "INPUT_PROXY_PORT"},
			Value:   "22",
		},
		&cli.DurationFlag{
			Name:    "proxy.timeout",
			Usage:   "proxy connection timeout",
			EnvVars: []string{"PLUGIN_PROXY_TIMEOUT", "PROXY_SSH_TIMEOUT", "INPUT_PROXY_TIMEOUT"},
		},
		&cli.StringSliceFlag{
			Name:    "envs",
			Usage:   "pass environment variable to shell script",
			EnvVars: []string{"PLUGIN_ENVS", "INPUT_ENVS"},
		},
		&cli.BoolFlag{
			Name:    "debug",
			Usage:   "debug mode",
			EnvVars: []string{"PLUGIN_DEBUG", "DEBUG", "INPUT_DEBUG"},
		},
	}

	// Override a template
	cli.AppHelpTemplate = `
________                                         _________ _________ ___ ___
\______ \_______  ____   ____   ____            /   _____//   _____//   |   \
 |    |  \_  __ \/  _ \ /    \_/ __ \   ______  \_____  \ \_____  \/    ~    \
 |    |   \  | \(  <_> )   |  \  ___/  /_____/  /        \/        \    Y    /
/_______  /__|   \____/|___|  /\___  >         /_______  /_______  /\___|_  /
        \/                  \/     \/                  \/        \/       \/
                                                    version: {{.Version}}
NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
COMMANDS:
{{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}{{if .Copyright }}
COPYRIGHT:
   {{.Copyright}}
   {{end}}{{if .Version}}
VERSION:
   {{.Version}}
   {{end}}
REPOSITORY:
    Github: https://github.com/appleboy/drone-ssh
`

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	scripts := c.StringSlice("script")
	if s := c.String("script.string"); s != "" {
		scripts = append(scripts, s)
	}
	plugin := Plugin{
		Config: Config{
			Key:            c.String("ssh-key"),
			KeyPath:        c.String("key-path"),
			Username:       c.String("user"),
			Password:       c.String("password"),
			Passphrase:     c.String("ssh-passphrase"),
			Host:           c.StringSlice("host"),
			Port:           c.Int("port"),
			Timeout:        c.Duration("timeout"),
			CommandTimeout: c.Duration("command.timeout"),
			Script:         scripts,
			ScriptStop:     c.Bool("script.stop"),
			Envs:           c.StringSlice("envs"),
			Debug:          c.Bool("debug"),
			Sync:           c.Bool("sync"),
			Proxy: easyssh.DefaultConfig{
				Key:        c.String("proxy.ssh-key"),
				KeyPath:    c.String("proxy.key-path"),
				User:       c.String("proxy.username"),
				Password:   c.String("proxy.password"),
				Passphrase: c.String("proxy.ssh-passphrase"),
				Server:     c.String("proxy.host"),
				Port:       c.String("proxy.port"),
				Timeout:    c.Duration("proxy.timeout"),
			},
		},
		Writer: os.Stdout,
	}

	return plugin.Exec()
}
