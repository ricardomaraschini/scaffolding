package main

import (
	"github.com/spf13/cobra"

	"app/ctrls"
	"app/infra/starter"
)

func init() {
	serve.Flags().String("http-address", ":8080", "listen tcp address for http server")
	serve.Flags().String("metrics-address", ":8181", "listen tcp address for metrics listen")
}

var serve = &cobra.Command{
	Use:   "serve",
	Short: "Starts application controllers",
	RunE: func(c *cobra.Command, args []string) error {
		haddr, err := c.Flags().GetString("http-address")
		if err != nil {
			return err
		}

		maddr, err := c.Flags().GetString("metrics-address")
		if err != nil {
			return err
		}

		starter := starter.NewEngine(
			ctrls.NewPromCtrl(ctrls.WithPromCtrlBindAddress(maddr)),
			ctrls.NewAppCtrl(ctrls.WithAppCtrlBindAddress(haddr)),
		)
		return starter.Start(c.Context())
	},
}
