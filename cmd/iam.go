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
	"strings"

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

		var maxUsernameLen int
		for _, user := range result.Users {
			if len(*user.UserName) > maxUsernameLen {
				maxUsernameLen = len(*user.UserName)
			}
		}
		var columnWidth int = 15
		if maxUsernameLen > columnWidth {
			columnWidth = maxUsernameLen
		}
		fmt.Printf("USERNAME%sLAST LOGIN MFA\n", strings.Repeat(" ", columnWidth-len("USERNAME")))
		fmt.Printf("--------%s----------\n", strings.Repeat(" ", columnWidth-len("USERNAME")))

   for _, user := range result.Users {
        if user.PasswordLastUsed == nil {
            fmt.Printf("%s%s%s",*user.UserName,strings.Repeat(" ", columnWidth-len(*user.UserName)),"NEVER")
        } else {
            fmt.Printf("%s%s%s",*user.UserName,strings.Repeat(" ", columnWidth-len(*user.UserName)),user.PasswordLastUsed)
        }
        // check if the user has MFA enabled
        mfa, _ := svc.ListMFADevices(&iam.ListMFADevicesInput{UserName: user.UserName})
        if len(mfa.MFADevices) > 0 {
            fmt.Printf("\tYES\n")
        } else {
		fmt.Printf("\tNO\n")
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
