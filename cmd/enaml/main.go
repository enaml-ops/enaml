package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/xchapter7x/enaml/generators"
	"github.com/xchapter7x/enaml/pull"
	"github.com/xchapter7x/enaml/run"
)

func init() {
	if c := os.Getenv("ENAML_CACHE_DIR"); c != "" {
		cacheDir = c
	}

	if o := os.Getenv("ENAML_OUTPUT_DIR"); o != "" {
		outputDir = o
	}

	os.MkdirAll(cacheDir, 0755)
}

const (
	CacheDir  = ".cache"
	OutputDir = "./enaml-gen"
)

var (
	Version   string
	cacheDir  = CacheDir
	outputDir = OutputDir
)

func main() {
	app := cli.NewApp()
	app.Version = Version
	app.Commands = []cli.Command{
		{
			Name:        "generate-jobs",
			Aliases:     []string{"gj"},
			Usage:       "generate-jobs <releaseurl>",
			Description: "generate golang structs for the jobs in a given release",
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
				d := run.NewDiffCmd(pull.Release{CacheDir: cacheDir}, c.Args()[0], c.Args()[1])
				err := d.All(os.Stdout)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			},
		},
		{
			Name:        "diff-job",
			Aliases:     []string{"dj"},
			Usage:       "diff-job <jobname> <releaseurl-A> <releaseurl-B>",
			Description: "show diff between jobs across 2 releases",
			Action: func(c *cli.Context) {
				d := run.NewDiffCmd(pull.Release{CacheDir: cacheDir}, c.Args()[1], c.Args()[2])
				err := d.Job(c.Args()[0], os.Stdout)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			},
		},
		{
			Name:        "show",
			Aliases:     []string{"sh"},
			Usage:       "show <releaseurl>",
			Description: "show all jobs and properties from the specified release",
			Action: func(c *cli.Context) {
				releaseRepo := pull.Release{CacheDir: cacheDir}
				s := run.NewShowCmd(releaseRepo, c.Args()[0])
				err := s.All(os.Stdout)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			},
		},
	}
	app.Run(os.Args)
}
