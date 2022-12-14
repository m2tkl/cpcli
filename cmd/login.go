/*
Copyright © 2022 m2tkl
*/
package cmd

import (
	"fmt"
	"syscall"

	"golang.org/x/term"

	"github.com/m2tkl/cpcli/cmd/internal"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login and store current session data",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		login()
	},
}

func init() {
	acCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func login() {
	client, _ := internal.NewClient("https://atcoder.jp", nil)

	if client.IsLoggedIn() {
		fmt.Println("Already logged in.")
		return
	}

	// Get username from stdin
	var username string
	fmt.Print("Enter username: ")
	fmt.Scan(&username)

	// Get password from stdin
	fmt.Print("Enter password: ")
	pwd, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("")

	client.Login(username, string(pwd))
}
