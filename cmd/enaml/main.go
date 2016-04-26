package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/xchapter7x/enaml/diff"
	"github.com/xchapter7x/enaml/generators"
	"github.com/xchapter7x/lo"
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
				diff := diff.NewDiff(cacheDir)
				diffset, err := diff.ReleaseDiff(c.Args()[0], c.Args()[1])
				displayDiffAndExit(diffset, err)
			},
		},
		{
			Name:        "diff-job",
			Aliases:     []string{"dj"},
			Usage:       "diff-jobs <jobname> <releaseurl-A> <releaseurl-B>",
			Description: "show diff between jobs across 2 releases",
			Action: func(c *cli.Context) {
				diff := diff.NewDiff(cacheDir)
				diffset, err := diff.JobDiffBetweenReleases(c.Args()[0], c.Args()[1], c.Args()[2])
				displayDiffAndExit(diffset, err)
			},
		},
	}
	app.Run(os.Args)
}

func displayDiffAndExit(diffset []string, err error) {
	if err == nil && len(diffset) == 0 {
		fmt.Println("no differences found")

	} else if len(diffset) > 0 {

		for i, v := range diffset {
			fmt.Println(i, ": ", v)
		}
		os.Exit(1)

	} else {
		lo.G.Error("error: ", err)
		os.Exit(2)
	}
	os.Exit(0)
}
