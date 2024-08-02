/*
Copyright © 2024 Muhammed Hussein Karimi info@karimi.dev

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"github.com/mhkarimi1383/health-checker/pkg/checkers"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/automaxprocs/maxprocs"
	"os"
	"time"
)

var (
	debug              bool
	trace              bool
	interval           time.Duration
	jsonLog            bool
	checkersConfigFile string
	checkConfigs       checkers.CheckConfigs
	chs                checkers.Checkers
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "health-checker",
	Short: "Health check your dependencies with ease",
	Long:  `Could be run in http server mode (e.g. K8s Sidecar) or check and exit mode (e.g. Init Container/Job)`,
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		if !jsonLog {
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		}
		defer maxprocs.Set(maxprocs.Logger(log.Printf))
		defer func() {
			viper.SetConfigFile(checkersConfigFile)
			if err := viper.ReadInConfig(); err != nil {
				log.Panic().Err(err).Msg("Unable to read checkers config")
			}
			if err := viper.Unmarshal(&checkConfigs); err != nil {
				log.Panic().Err(err).Msg("Unable to parse checkers config")
			}
			var err error
			if chs, err = checkers.ConfigsToCheckers(checkConfigs); err != nil {
				log.Panic().Err(err).Msg("Unable to parse configs")
			}
		}()
		if trace {
			zerolog.SetGlobalLevel(zerolog.TraceLevel)
			return
		}
		if debug {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
			return
		}
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		maxprocs.Set(maxprocs.Logger(log.Printf))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable Debug mode")
	rootCmd.PersistentFlags().BoolVarP(&trace, "trace", "t", false, "Enable Trace mode (Also enables debug mode)")
	rootCmd.PersistentFlags().BoolVarP(&jsonLog, "json-log", "j", false, "Enable logging in json format")
	rootCmd.PersistentFlags().StringVarP(&checkersConfigFile, "checkers", "c", "checkers.yaml", "Configuration file (Accepting JSON, TOML, YAML, HCL, INI, envfile or Java properties formats)")
	rootCmd.PersistentFlags().DurationVarP(&interval, "interval", "i", 2*time.Second, "Interval between each check")
}
