/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/m2tkl/cpcli/cmd/internal"
	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("setup called")
		setup("abc064")
	},
}

func init() {
	acCmd.AddCommand(setupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func setup(contest string) {
	acConfig := internal.NewAcConfig()
	c, _ := internal.NewClient(acConfig.Endpoint, nil)

	tasks := c.FetchContestTasks(contest)

	path := filepath.Join(acConfig.Dir + "/config.json")
	jsonText, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var config Config
	json.Unmarshal([]byte(jsonText), &config)
	fmt.Println(config.Dir)

	contestDirPath := config.Dir + "/contests/" + contest

	err = os.MkdirAll(contestDirPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range tasks {
		fmt.Println(k)
		fmt.Println(v)
		taskCases := c.FetchSampleTestCases(v)

		taskDirPath := contestDirPath + "/" + k

		// Create task dirs
		err = os.MkdirAll(taskDirPath+"/tests", os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}

		// Create main.go
		// TODO: copy template (option)
		if !exists(taskDirPath + "/main.go") {
			fp, err := os.Create(taskDirPath + "/main.go")
			if err != nil {
				log.Fatal(err)
				return
			}
			defer fp.Close()
		}

		var taskTestDir string
		for i, taskCase := range taskCases {
			fmt.Println(taskCase)
			taskTestDir = taskDirPath + "/tests/" + strconv.Itoa(i+1)
			err := os.Mkdir(taskTestDir, os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}

			writeLines(taskTestDir+"/in.txt", taskCase.In)
			writeLines(taskTestDir+"/out.txt", taskCase.Out)
		}
	}
}

func writeLines(filePath string, value string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	fmt.Fprintln(f, value)
	return nil
}
