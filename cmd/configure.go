/*
Copyright © 2022 m2tkl
*/
package cmd

import (
	"encoding/json"
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

	cwd, _ := os.Getwd()

	// NOTE:
	// 		If config file already exists,
	// 		confirm whether to overwrite it while displaying the contents.

	// TODO: Receive stdin

	config := Config{
		Dir:      cwd,
		Template: "",
	}

	file, _ := json.MarshalIndent(config, "", " ")
	_ = ioutil.WriteFile(acConfig.Dir+"/config.json", file, 0644)

}
