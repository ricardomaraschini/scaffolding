package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

func main() {
	sigs := []os.Signal{syscall.SIGTERM, syscall.SIGINT}
	ctx, cancel := signal.NotifyContext(context.Background(), sigs...)
	defer cancel()

	root := &cobra.Command{
		Use:          "app",
		SilenceUsage: true,
	}
	root.AddCommand(version, serve)
	if err := root.ExecuteContext(ctx); err != nil {
		os.Exit(2)
	}
}
