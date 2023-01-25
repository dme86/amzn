package bucket

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/cobra"
)

const (
    YYYYMMDD = "2006-01-02"
)

var RootCmd = &cobra.Command{
	Use:   "amzn",
	Short: "Minimal aws CLI",
}

func init() {
	RootCmd.AddCommand(listBucketsCmd)
}

var listBucketsCmd = &cobra.Command{
	Use:   "list-buckets",
	Short: "Lists the names and creation dates of all S3 buckets",
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

		fmt.Printf("%-"+fmt.Sprintf("%d", longestBucketName)+"s %-20s\n", "BUCKET NAME", "CREATION DATE")
		fmt.Printf("%-"+fmt.Sprintf("%d", longestBucketName)+"s %-20s\n", "-----------", "-------------")
		for _, b := range result.Buckets {
			fmt.Printf("%*s %s\n", -longestBucketName, aws.StringValue(b.Name), b.CreationDate.Format(YYYYMMDD))
		}
	},
}

