package cmd

import (
	"go_producer_consumer/internal/cli"
	"os"

	"github.com/spf13/cobra"
)

var (
	role              string
	msgCount          int
	consumerGroupName string
)

var rootCommand = &cobra.Command{
	Use:   "go_producer_consumer",
	Short: "Kafka Producer/Consumer for multiple kafka versions",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cli.Run(role, msgCount, consumerGroupName); err != nil {
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCommand.Flags().StringVar(&role, "role", "", "Role to run as. available roles are producer, consumer, both")
	rootCommand.Flags().IntVar(&msgCount, "msg-count", 10, "Number of messages to produce per topic")
	rootCommand.Flags().StringVar(&consumerGroupName, "consumer-group", "juniper-group", "Consumer Group name")
	rootCommand.MarkFlagRequired("role")
}
