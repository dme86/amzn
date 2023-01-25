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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/cobra"
)

const (
	YYYYMMDD = "2006-01-02"
	HHMMSS   = "12:05:03"
)

var noHeader bool

// s3Cmd represents the s3 command
var s3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "List of S3 buckets with creation date",
	Long:  `Foo`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.HelpFunc()(cmd, args)
	},
}

var s3lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List S3 buckets",
	Run: func(cmd *cobra.Command, args []string) {
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

		svc := s3.New(sess)

		result, err := svc.ListBuckets(nil)
		if err != nil {
			fmt.Println("Failed to list buckets", err)
			return
		}

		longestBucketName := 0
		for _, b := range result.Buckets {
			if len(aws.StringValue(b.Name)) > longestBucketName {
				longestBucketName = len(aws.StringValue(b.Name))
			}
		}
		if !noHeader {
			fmt.Printf("%-"+fmt.Sprintf("%d", longestBucketName)+"s %-20s\n", "BUCKET NAME", "CREATION DATE")
			fmt.Printf("%-"+fmt.Sprintf("%d", longestBucketName)+"s %-20s\n", "-----------", "-------------")
		}
		for _, b := range result.Buckets {
			fmt.Printf("%*s %s %s\n", -longestBucketName, aws.StringValue(b.Name), b.CreationDate.Format(YYYYMMDD), b.CreationDate.Format("15:04:05"))
		}
	},
}

func init() {
	rootCmd.AddCommand(s3Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// s3Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	s3Cmd.Flags().BoolVarP(&noHeader, "no-header", "", false, "Do not print column headers")
	s3Cmd.AddCommand(s3lsCmd)
}
