package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
)

func main() {
	if err := NewCLI().RootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}

var version string

type CLI struct {
	doRun func() error
}

func NewCLI() *CLI {
	return &CLI{
		doRun: doRun,
	}
}

func (c *CLI) RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "ms1-server",
		Version:           version,
		Short:             "MS1 Server",
		Long:              "MS1 Server",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		SilenceErrors:     true,
		SilenceUsage:      true,
		Args:              cobra.RangeArgs(0, 0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	var (
		b      bytes.Buffer
		config Config
	)

	if err := envconfig.Usagef("MS1", &config, &b, envconfig.DefaultTableFormat); err != nil {
		panic(err)
	}

	cmd.SetUsageTemplate(cmd.UsageTemplate() + "Environment variables\n\n" + b.String())

	cmd.SetHelpCommand(&cobra.Command{Hidden: true})
	cmd.AddCommand(c.runCmd())

	return cmd
}

func (c *CLI) runCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Run the server.",
		Long:  "Run the server.",
		Args:  cobra.RangeArgs(0, 0),
		RunE: func(_ *cobra.Command, _ []string) error {
			return c.doRun()
		},
	}
}

func doRun() error {
	config, err := NewConfigFromEnv()
	if err != nil {
		return err
	}

	logger := NewLogger(config.LogLevel)

	grpcServer := &GrpcServer{
		Addr:              fmt.Sprintf(":%d", config.GrpcPort),
		EchoServiceServer: NewEchoServiceServer(),
	}

	server := &Server{
		Logger:     logger,
		GrpcServer: grpcServer,
	}

	return server.Run()
}
