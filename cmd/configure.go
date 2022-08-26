/*
Copyright © 2022 m2tkl
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/m2tkl/cpcli/cmd/internal"
	"github.com/spf13/cobra"
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure 'ac' command setting",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		configure()
	},
}

func init() {
	acCmd.AddCommand(configureCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configureCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configureCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type Config struct {
	Dir      string `json:"dir"`
	Template string `json:"template"`
}

func configure() {
	acConfig := internal.NewAcConfig()

	err := os.MkdirAll(acConfig.Dir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	// NOTE:
	// 		If config file already exists,
	// 		confirm whether to overwrite it while displaying the contents.

	configJson, err := ioutil.ReadFile(acConfig.Dir + "/config.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// TODO: impl '--reset' option

	var config Config
	json.Unmarshal([]byte(configJson), &config)
	fmt.Println(config.Dir)

	var dirPath string

	if config.Dir != "" {
		fmt.Printf("(%s) Enter directory path: ", config.Dir)
		fmt.Scanln(&dirPath)

		// If not entered, use saved directory path
		if dirPath == "" {
			dirPath = config.Dir
		}

	} else {
		fmt.Print("Enter directory path: ")
		fmt.Scanln(&dirPath)

		// If not entered, use current directory
		if dirPath == "" {
			dirPath, _ = os.Getwd()
		}
	}

	var templatePath string

	if config.Template != "" {
		fmt.Printf("(%s) Enter template path: ", config.Template)
		fmt.Scanln(&templatePath)

		// If not entered, use saved directory path
		if templatePath == "" {
			templatePath = config.Template
		}

	} else {
		fmt.Print("Enter template path: ")
		fmt.Scanln(&templatePath)

		// TODO: use default value
		// If not entered, not set path
		if templatePath == "" {
			templatePath = ""
		}
	}

	newConfig := Config{
		Dir:      dirPath,
		Template: templatePath,
	}

	file, _ := json.MarshalIndent(newConfig, "", " ")
	_ = ioutil.WriteFile(acConfig.Dir+"/config.json", file, 0644)
}
