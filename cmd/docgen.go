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
	"io/fs"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// docGenCmd represents the docGen command
var docGenCmd = &cobra.Command{
	Use:   "docGen",
	Short: "Generate Cli docs",
	Long:  `Generates Cli docs in docs/ directory`,
	Run: func(cmd *cobra.Command, args []string) {
		err := os.MkdirAll("docs/", fs.ModeDir|fs.ModePerm)
		if err != nil {
			log.Panic().Err(err).Msg("docs/ Dir creation error")
		}
		err = doc.GenMarkdownTree(rootCmd, "docs/")
		if err != nil {
			log.Panic().Err(err).Msg("Markdown DocGen Error")
		}

		err = os.MkdirAll("docs/man/", fs.ModeDir|fs.ModePerm)
		if err != nil {
			log.Panic().Err(err).Msg("docs/man/ Dir creation error")
		}
		err = doc.GenManTree(rootCmd, &doc.GenManHeader{
			Title:   "HEALTH CHECKER",
			Section: "3",
		}, "docs/man")
		if err != nil {
			log.Panic().Err(err).Msg("Man DocGen Error")
		}

		err = os.MkdirAll("docs/rst/", fs.ModeDir|fs.ModePerm)
		if err != nil {
			log.Panic().Err(err).Msg("docs/rst/ Dir creation error")
		}
		err = doc.GenReSTTree(rootCmd, "docs/rst/")
		if err != nil {
			log.Panic().Err(err).Msg("RST DocGen Error")
		}

		err = os.MkdirAll("docs/yaml/", fs.ModeDir|fs.ModePerm)
		if err != nil {
			log.Panic().Err(err).Msg("docs/yaml/ Dir creation error")
		}
		err = doc.GenYamlTree(rootCmd, "docs/yaml/")
		if err != nil {
			log.Panic().Err(err).Msg("YAML DocGen Error")
		}
	},
}

func init() {
	rootCmd.AddCommand(docGenCmd)
}
