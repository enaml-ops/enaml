package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/xchapter7x/enaml/generators"
)

func init() {
	if c := os.Getenv("CACHE_DIR"); c != "" {
		cacheDir = c
	}
	os.MkdirAll(cacheDir, 0755)
}

const (
	CacheDir  = ".cache"
	OutputDir = "./releasejobs"
)

var (
	cacheDir = CacheDir
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:    "generate-jobs",
			Aliases: []string{"gj"},
			Usage:   "generate golang structs for the jobs in a given release",
			Action: func(c *cli.Context) {
				generators.GenerateReleaseJobsPackage(c.Args().First(), cacheDir, OutputDir)
				println("completed generating release job structs for ", c.Args().First())
			},
		},
		{
			Name:    "diff-release",
			Aliases: []string{"dr"},
			Usage:   "show a diff between 2 releases given",
			Action: func(c *cli.Context) {
				println("unimplemented")
				println("release job properties diff", c.Args().First())
			},
		},
		{
			Name:    "diff-job",
			Aliases: []string{"dj"},
			Usage:   "show diff between jobs across 2 releases",
			Action: func(c *cli.Context) {
				println("unimplemented")
				println("release job properties diff", c.Args().First())
			},
		},
	}
	app.Run(os.Args)
}
