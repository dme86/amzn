/*
Copyright © 2023 Daniel <dme86> Meier

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	//		"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

// iamCmd represents the iam command
var iamCmd = &cobra.Command{
	Use:   "iam",
	Short: "A brief description of your command",
	Long:  `Foo`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.HelpFunc()(cmd, args)
	},
}

var iamlsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List iam users",
	Run: func(cmd *cobra.Command, args []string) {
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

		svc := iam.New(sess)

		result, err := svc.ListUsers(&iam.ListUsersInput{})
		if err != nil {
			fmt.Println("Error", err)
			os.Exit(1)
		}
		// Print the column headers
		fmt.Printf("%-20s %-20s %-5s\n", "USERNAME", "LAST LOGIN", "MFA")

		// Loop through the list of users and print their details
		for _, user := range result.Users {
			// Get the user's login profile
			userProfile, err := svc.GetUser(&iam.GetUserInput{
				UserName: user.UserName,
			})
			if err != nil {
				fmt.Println("Error", err)
				continue
			}

			// Get the user's MFA status
			mfa, err := svc.ListMFADevices(&iam.ListMFADevicesInput{
				UserName: user.UserName,
			})
			if err != nil {
				fmt.Println("Error", err)
				continue
			}

			// Print the user's details
			if userProfile.User.PasswordLastUsed != nil {
				fmt.Printf("%-20s %-20s", *user.UserName, userProfile.User.PasswordLastUsed.Format(time.RFC1123))
			} else {
				fmt.Printf("%-20s %-20s", *user.UserName, "NEVER")
			}
			if len(mfa.MFADevices) > 0 {
				fmt.Printf("%-5s\n", "YES")
			} else {
				fmt.Printf("%-5s\n", "NO")
			}

		}
	},
}

func init() {
	rootCmd.AddCommand(iamCmd)
	iamCmd.AddCommand(iamlsCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// iamCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// iamCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
