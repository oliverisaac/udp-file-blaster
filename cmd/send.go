/*
Copyright © 2024 Oliver Isaac

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
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var destAddress string
var destPort int
var absolutePaths bool

// sendCmd represents the receive command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Start a client to send files to a receiver",
	Long:  `Once you've started a receiver, use this to send files to the receiver`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("not impl")
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)

	sendCmd.Flags().StringVarP(&destAddress, "address", "a", "", "Address to send to")
	sendCmd.MarkFlagRequired("address")

	sendCmd.Flags().IntVarP(&destPort, "port", "p", 9876, "Port to send to")
	sendCmd.Flags().BoolVar(&absolutePaths, "absolute", false, "Strip directory paths to and just use base filenames")
}
