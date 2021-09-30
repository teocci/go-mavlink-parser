// Package cmd
// Created by RTT.
// Author: teocci@yandex.com on 2021-Sep-27
package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/teocci/go-mavlink-parser/src/cmd/cmdapp"
	"github.com/teocci/go-mavlink-parser/src/config"
	"github.com/teocci/go-mavlink-parser/src/core"
	"github.com/teocci/go-mavlink-parser/src/data"
	"github.com/teocci/go-mavlink-parser/src/logger"
)

var (
	app = &cobra.Command{
		Use:           cmdapp.Name,
		Short:         cmdapp.Short,
		Long:          cmdapp.Long,
		PreRunE:       validate,
		RunE:          runE,
		SilenceErrors: false,
		SilenceUsage:  true,
	}

	host      string
	port      int64
	connID    int64
	moduleTag string
	droneID   int64
	flightID  int64
)

// Add supported cli commands/flags
func init() {
	cobra.OnInitialize(initConfig)

	app.Flags().StringVarP(&host, cmdapp.HName, cmdapp.HShort, host, cmdapp.HDesc)
	app.Flags().StringVarP(&moduleTag, cmdapp.MName, cmdapp.MShort, moduleTag, cmdapp.MDesc)

	app.Flags().Int64VarP(&port, cmdapp.PName, cmdapp.PShort, port, cmdapp.PDesc)
	app.Flags().Int64VarP(&connID, cmdapp.CName, cmdapp.CShort, connID, cmdapp.CDesc)
	app.Flags().Int64VarP(&droneID, cmdapp.DName, cmdapp.DShort, droneID, cmdapp.DDesc)
	app.Flags().Int64VarP(&flightID, cmdapp.FName, cmdapp.FShort, flightID, cmdapp.FDesc)

	_ = app.MarkFlagRequired(cmdapp.HName)
	_ = app.MarkFlagRequired(cmdapp.MName)

	_ = app.MarkFlagRequired(cmdapp.PName)
	_ = app.MarkFlagRequired(cmdapp.CName)
	_ = app.MarkFlagRequired(cmdapp.DName)
	_ = app.MarkFlagRequired(cmdapp.FName)

	config.AddFlags(app)
}

// Load config
func initConfig() {
	if err := config.LoadConfigFile(); err != nil {
		log.Fatal(err)
	}

	config.LoadLogConfig()
}

func validate(ccmd *cobra.Command, args []string) error {
	if config.Version {
		fmt.Printf(cmdapp.VersionTemplate, cmdapp.Name, cmdapp.Version, cmdapp.Commit)

		return nil
	}

	if !config.Verbose {
		ccmd.HelpFunc()(ccmd, args)

		return fmt.Errorf("")
	}

	return nil
}

func runE(ccmd *cobra.Command, args []string) error {
	var err error
	config.Log, err = logger.New(config.LogConfig)
	if err != nil {
		return ErrCanNotLoadLogger(err)
	}

	initData := data.InitConf{
		Host:      host,
		Port:      port,
		ConnID:    connID,
		ModuleTag: moduleTag,
		DroneID:   droneID,
		FlightID:  flightID,
	}
	// Make channel for errors
	errs := make(chan error)
	go func() {
		errs <- core.Start(initData)
	}()

	// Break if any of them return an error (blocks exit)
	if err := <-errs; err != nil {
		config.Log.Fatal(err)
	}

	return err
}

func Execute() {
	err := app.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
