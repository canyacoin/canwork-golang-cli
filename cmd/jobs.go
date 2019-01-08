package cmd

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/spf13/cobra"
	"gitlab.com/canya-com/shared/data-structures/canwork"
)

var (
	JobID string
)

var jobsCmd = &cobra.Command{
	Use:   "jobs",
	Short: "get information about jobs",
	Long:  `Jobs collects information for a given job, or lists recent jobs.`,
	Run: func(cmd *cobra.Command, args []string) {

		fireStore = newFirestoreClient(KeyFile)
		ctx := context.Background()

		if JobID != "" {
			fmt.Printf("+ getting job information for jobID: %s\n", JobID)

			js, err := fireStore.Collection("jobs").Doc(JobID).Get(ctx)
			if err != nil {
				panic(err)
			}
			var job canwork.JobDocument
			js.DataTo(&job)
			if job.ID == "" {
				fmt.Printf("job ID: %s could not be found\n\n", JobID)
			}
			fmt.Printf("+ found job with title : %s (%s)\n", job.Information.Title, JobID)

			cs, err := fireStore.Collection("users").Doc(job.ClientID).Get(ctx)
			if err != nil {
				panic(err)
			}
			var client canwork.UserDocument
			cs.DataTo(&client)
			fmt.Printf("+ client               : %s (%s)\n", client.Name, job.ClientID)

			ps, err := fireStore.Collection("users").Doc(job.ProviderID).Get(ctx)
			if err != nil {
				panic(err)
			}
			var provider canwork.UserDocument
			ps.DataTo(&provider)
			fmt.Printf("+ provider             : %s (%s)\n", provider.Name, job.ProviderID)

			// Get job balance

			// Get escrow balance

			// Get actions

			fmt.Print("\n")

		} else {
			fmt.Printf("+ getting job information for recent jobs\n")
			ji := fireStore.Collection("jobs").OrderBy("deadline", firestore.Asc).Documents(ctx)
			for {
				d, err := ji.Next()
				if err != nil {
					break
				}
				var job canwork.JobDocument
				d.DataTo(&job)
				if job.ID == "" {
					continue
				}
				fmt.Printf("job id: %s\n", job.ID)
				fmt.Printf("title : %s\n\n", job.Information.Title)
			}
		}

		// u, err := fireBase.GetUserByEmail(ctx, Email)
		// if err != nil {
		// 	logger.Fatalf("error getting user %s: %v\n", Email, err)
		// }

		// token, err := fireBase.CustomToken(ctx, u.UID)
		// if err != nil {
		// 	logger.Fatalf("error minting new token for userID %s: %v\n", u.UID, err)
		// }

		// fmt.Printf("+ json web token copied to clipboard\n")
		// clipboard.WriteAll(token)

		// if Verbose {
		// 	fmt.Println(token)
		// }

	},
}

func init() {
	rootCmd.AddCommand(jobsCmd)

	jobsCmd.Flags().StringVarP(&KeyFile, "keyFile", "k", "", "path to the json firebase credentials key file (required)")
	jobsCmd.MarkFlagRequired("keyFile")

	jobsCmd.Flags().StringVarP(&JobID, "jobID", "j", "", "job id to examine")

	// loginCmd.Flags().StringVarP(&Email, "email", "e", "", "email address of the user to login as (required)")
	// loginCmd.MarkFlagRequired("email")

	jobsCmd.Flags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
