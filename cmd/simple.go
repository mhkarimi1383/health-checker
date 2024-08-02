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
	"time"

	"github.com/mhkarimi1383/health-checker/pkg/checkers"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// simpleCmd represents the simple command
var simpleCmd = &cobra.Command{
	Use:   "simple",
	Short: "Just do a health check and exit or wait to service to become available",
	Long: `Will do the health check and exit or wait to service to become available,

also with correct exit code based on status`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("Starting in simple mode")
		success := false
		reRun := true
		for reRun {
			status := checkers.RunChecks(chs)
			success = true
			for name, s := range status {
				log.Info().Dur("latency", s.Latency).Bool("isAlive", s.IsAlive).Any("error", s.Error).Str("name", name).Str("type", s.Type).Msg("Health Check status")
				if !s.IsAlive {
					success = false
				}
			}
			reRun = (!success && wait)
			if reRun {
				time.Sleep(interval)
			}
		}
		if !success {
			log.Fatal().Msg("One or more health checks failed")
		}
		log.Info().Msg("All checks are passed")
	},
}

var (
	wait bool
)

func init() {
	rootCmd.AddCommand(simpleCmd)
	simpleCmd.Flags().BoolVarP(&wait, "wait", "w", false, "Wait for services (do not exit immediately)")
}
