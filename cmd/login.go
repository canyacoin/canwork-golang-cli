package cmd

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login to firebase via email",
	Long: `Login authenticates a user and mint's a JWT for the user you are logging in as.
This is useful for calling API's that require a valid (not expired) Json Web Token:

The returned token will be copied into the clipboard for your convenience.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("+ minting json web token for user: %s\n", Email)

		// fireStore = newFirestoreClient(KeyFile)
		fireBase = newFirebaseClient(KeyFile)

		u, err := fireBase.GetUserByEmail(ctx, Email)
		if err != nil {
			logger.Fatalf("error getting user %s: %v\n", Email, err)
		}

		token, err := fireBase.CustomToken(ctx, u.UID)
		if err != nil {
			logger.Fatalf("error minting new token for userID %s: %v\n", u.UID, err)
		}

		fmt.Printf("+ json web token copied to clipboard\n")
		clipboard.WriteAll(token)

		if Verbose {
			fmt.Println(token)
		}

	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVarP(&KeyFile, "keyFile", "k", "", "path to the json firebase credentials key file (required)")
	loginCmd.MarkFlagRequired("keyFile")

	loginCmd.Flags().StringVarP(&Email, "email", "e", "", "email address of the user to login as (required)")
	loginCmd.MarkFlagRequired("email")

	loginCmd.Flags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
