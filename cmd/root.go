package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	homedir "github.com/mitchellh/go-homedir"
	logging "github.com/op/go-logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

var (
	logger    = logging.MustGetLogger("main")
	fireStore *firestore.Client
	fireBase  *auth.Client
	cfgFile   string
	ctx       = context.Background()
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "canwork-cli",
	Short: "Command line utilities for CANWork",
	Long:  `A set of utility functions designed to improve the quality of life of developers and support technicians alike.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	logging.SetLevel(logging.DEBUG, "main")
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".canwork-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".canwork-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func newFirebaseClient(keyFile string) *auth.Client {
	sa := option.WithCredentialsFile(keyFile)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		logger.Fatalf("unable parse service account credentials: %v", err)
	}

	fbc, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}
	return fbc
}

func newFirestoreClient(keyFile string) *firestore.Client {
	sa := option.WithCredentialsFile(keyFile)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		logger.Fatalf("unable parse service account credentials: %v", err)
	}
	fsc, err := app.Firestore(ctx)
	if err != nil {
		logger.Fatalf("unable to establish firestore connection: %v", err)
	}
	return fsc
}
