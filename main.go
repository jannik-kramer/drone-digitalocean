package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
	"os"
)

// build number set at compile-time
var build = "0"

// Version set at compile-time
var Version string

var plugin = Plugin{
	Writer: os.Stdout,
}

func main() {
	if Version == "" {
		Version = fmt.Sprintf("0.1.0+%s", build)
	}

	// Load env-file if it exists first
	if filename, found := os.LookupEnv("PLUGIN_ENV_FILE"); found {
		_ = godotenv.Load(filename)
	}

	app := cli.NewApp()
	app.Name = "Drone DigitalOcean"
	app.Usage = "Deploy your application easly to DigitalOcean"
	app.Copyright = "Copyright (c) Jannik Kramer"
	app.Authors = []cli.Author{
		{
			Name:  "Jannik Kramer",
			Email: "mail@jannikkramer.de",
		},
	}
	app.Action = run
	app.Version = Version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "loadbalancer,b",
			Usage:  "loadbalancer to deploy to",
			EnvVar: "PLUGIN_LOADBALANCER",
			Destination: &plugin.Config.Loadbalancer,
		},
		cli.StringFlag{
			Name: "tag",
			Usage: "tag to deploy to",
			EnvVar: "PLUGIN_TAG,TAG",
			Destination: &plugin.Config.Tag,
		},
		cli.StringFlag{
			Name:   "user,u",
			Usage:  "connect as user",
			EnvVar: "PLUGIN_SSH_USER,SSH_USER",
			Destination: &plugin.Config.User,
		},
		cli.StringFlag{
			Name:   "key-path,i",
			Usage:  "ssh private key path",
			EnvVar: "PLUGIN_SSH_KEY_PATH,SSH_KEY_PATH",
			Destination: &plugin.Config.KeyPath,
		},
		cli.StringFlag{
			Name:   "key",
			Usage:  "ssh private key",
			EnvVar: "PLUGIN_SSH_KEY,SSH_KEY",
			Destination: &plugin.Config.Key,
		},
		cli.StringFlag{
			Name:   "password",
			Usage:  "ssh password",
			EnvVar: "PLUGIN_SSH_PASSWORD,SSH_PASSWORD",
			Destination: &plugin.Config.Password,
		},
		cli.IntFlag{
			Name:   "port,p",
			Usage:  "ssh port",
			EnvVar: "PLUGIN_SSH_PORT,SSH_PORT",
			Destination: &plugin.Config.Port,
		},
		cli.DurationFlag{
			Name:   "timeout",
			Usage:  "ssh timeout",
			EnvVar: "PLUGIN_SSH_TIMEOUT,SSH_TIMEOUT",
			Destination: &plugin.Config.Timeout,
		},
		cli.StringFlag{
			Name:   "pat,t",
			Usage:  "digitalocean personal access token",
			EnvVar: "PLUGIN_PAT,PAT",
			Destination: &plugin.Config.Pat,
		},
		cli.StringFlag{
			Name:   "source",
			Usage:  "folder to copy from",
			EnvVar: "PLUGIN_SOURCE,SOURCE",
			Destination: &plugin.Config.SourcePath,
		},
		cli.StringFlag{
			Name:   "target",
			Usage:  "folder to copy to",
			EnvVar: "PLUGIN_TARGET,TARGET",
			Destination: &plugin.Config.TargetPath,
		},
		cli.StringSliceFlag{
			Name:   "pre-sync",
			Usage:  "commands to run before sync",
			EnvVar: "PLUGIN_PRE_SYNC,PRE_SYNC",
		},
		cli.StringSliceFlag{
			Name:   "post-sync",
			Usage:  "commands to run after sync",
			EnvVar: "PLUGIN_POST_SYNC,POST_SYNC",
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

// this runs due to "app.Action = run" in main()
func run(c *cli.Context) error {
	plugin.Config.PreSync = c.StringSlice("pre-sync")
	plugin.Config.PostSync = c.StringSlice("post-sync")

	plugin.Config.KeyPath = "id_sha"

	if err := plugin.Exec(); err != nil {
		return err
	}

	return nil
}
