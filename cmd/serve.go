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
	"github.com/mhkarimi1383/health-checker/pkg/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Acts as http server",
	Long: `Will expose health checks

Accepts HEAD method if status is only needed
Also expose an API and HTML view of your health check

Could be used as a sidecar container/process to handle them`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Str("listenAdress", listenAddress).Msg("Starting in server mode")
		if err := server.Start(listenAddress, interval, chs); err != nil {
			log.Panic().Err(err).Msg("HTTP Server Error")
		}
	},
}

var (
	listenAddress string
)

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&listenAddress, "listen-address", "l", "127.0.0.1:2200", "Help message for toggle")
}
