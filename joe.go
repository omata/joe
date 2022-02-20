package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/urfave/cli"
)

const (
	gitignoreUrl = "https://github.com/github/gitignore/archive/main.zip"
	version      = "1.0.3"
	dataDir      = ".joe-data"
	joe          = `
 ▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄▄
▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌
 ▀▀▀▀▀█░█▀▀▀ ▐░█▀▀▀▀▀▀▀█░▌▐░█▀▀▀▀▀▀▀▀▀
      ▐░▌    ▐░▌       ▐░▌▐░▌
      ▐░▌    ▐░▌       ▐░▌▐░█▄▄▄▄▄▄▄▄▄
      ▐░▌    ▐░▌       ▐░▌▐░░░░░░░░░░░▌
      ▐░▌    ▐░▌       ▐░▌▐░█▀▀▀▀▀▀▀▀▀
      ▐░▌    ▐░▌       ▐░▌▐░▌
 ▄▄▄▄▄█░▌    ▐░█▄▄▄▄▄▄▄█░▌▐░█▄▄▄▄▄▄▄▄▄
▐░░░░░░░▌    ▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌
 ▀▀▀▀▀▀▀      ▀▀▀▀▀▀▀▀▀▀▀  ▀▀▀▀▀▀▀▀▀▀▀
`
)

var (
	errLogger   = log.New(os.Stderr, "", 0)
	userHome, _ = os.UserHomeDir()
	dataPath    = path.Join(userHome, dataDir)
)

func findGitignores() (a map[string]string, err error) {
	_, err = ioutil.ReadDir(dataPath)
	if err != nil {
		return nil, err
	}

	filelist := make(map[string]string)
	filepath.Walk(dataPath, func(filepath string, info os.FileInfo, err error) error {
		if strings.HasSuffix(info.Name(), ".gitignore") {
			name := strings.ToLower(strings.Replace(info.Name(), ".gitignore", "", 1))
			filelist[name] = filepath
		}
		return nil
	})
	return filelist, nil
}

func availableFiles() (a []string, err error) {
	gitignores, err := findGitignores()
	if err != nil {
		return nil, err
	}

	availableGitignores := []string{}
	for key, _ := range gitignores {
		availableGitignores = append(availableGitignores, key)
	}

	return availableGitignores, nil
}

func search(arg string) {
	gitignores, err := findGitignores()
	if err != nil {
		errLogger.Fatal(err)
	}

	for ig := range gitignores {
		b, err := regexp.MatchString(arg, ig)
		if err != nil {
			errLogger.Println(err)
		} else if b {
			fmt.Println(ig)
		}
	}
}

func generate(args string) {
	names := strings.Split(args, ",")

	gitignores, err := findGitignores()
	if err != nil {
		errLogger.Fatal(err)
	}

	notFound := []string{}
	output := "### Operating Systems ###\n# Darwin\n.DS_Store\n"
	for index, name := range names {
		if filepath, ok := gitignores[strings.ToLower(name)]; ok {
			bytes, err := ioutil.ReadFile(filepath)
			if err == nil {
				output += "\n#### " + name + " ####\n"
				output += string(bytes)
				if index < len(names)-1 {
					output += "\n"
				}
				continue
			}
		} else {
			notFound = append(notFound, name)
		}
	}

	if len(notFound) > 0 {
		errLogger.Printf("Unsupported files: %s\n", strings.Join(notFound, ", "))
		errLogger.Fatal("Run `joe ls` to see list of available gitignores.")
	}
	if len(output) > 0 {
		output = "#### joe made this: http://goel.io/joe\n" + output
	}

	fmt.Println(output)
}

func main() {
	app := cli.NewApp()
	app.Name = joe
	app.Usage = "generate .gitignore files from the command line"
	app.UsageText = "joe command [arguments...]"
	app.Version = version
	// app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		{
			Name:    "ls",
			Aliases: []string{"list"},
			Usage:   "list all available files",
			Action: func(c *cli.Context) error {
				availableGitignores, err := availableFiles()
				if err != nil {
					errLogger.Fatal(err)
					return err
				}
				fmt.Printf("%d supported .gitignore files:\n", len(availableGitignores))
				sort.Strings(availableGitignores)
				for _, gnore := range availableGitignores {
					fmt.Println(gnore)
				}
				return nil
			},
		},
		{
			Name:    "u",
			Aliases: []string{"update"},
			Usage:   "update all available gitignore files",
			Action: func(c *cli.Context) error {
				fmt.Println("Updating gitignore files..")
				err := RemoveContents(dataPath)
				if err != nil {
					errLogger.Fatal(err)
				}
				err = DownloadFiles(gitignoreUrl, dataPath)
				if err != nil {
					errLogger.Fatal(err)
					return err
				}
				return nil
			},
		},
		{
			Name:    "g",
			Aliases: []string{"generate"},
			Usage:   "generate gitignore files",
			Action: func(c *cli.Context) error {
				if c.NArg() != 1 {
					cli.ShowAppHelp(c)
				} else {
					generate(c.Args()[0])
				}
				return nil
			},
		},
		{
			Name:    "s",
			Aliases: []string{"search"},
			Usage:   "search for gitignore files (one word per query)",
			Action: func(c *cli.Context) error {
				if c.NArg() != 1 {
					cli.ShowAppHelp(c)
				} else {
					search(c.Args()[0])
				}
				return nil
			},
		},
	}
	app.Run(os.Args)
}
