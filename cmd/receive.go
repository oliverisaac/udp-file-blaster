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

var listenAddress string
var listenPort int
var baseDirectory string
var stripPaths bool

// receiveCmd represents the receive command
var receiveCmd = &cobra.Command{
	Use:   "receive",
	Short: "Start a client to receive files from a sender",
	Long:  `Starting a client will tell you where to send the files to`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("Not impl")
	},
}

func init() {
	rootCmd.AddCommand(receiveCmd)

	receiveCmd.Flags().StringVarP(&listenAddress, "address", "a", "0.0.0.0", "Address to listen on")
	receiveCmd.Flags().IntVarP(&listenPort, "port", "p", 9876, "Port to listen on")
	receiveCmd.Flags().StringVarP(&baseDirectory, "dir", "d", "/", "Base directory to place the received files. Files will be placed relative to this directory, so if you receive file '/var/logs/messages' and specify --dir=/example then the file will land in /example/var/logs/messages")
	receiveCmd.Flags().BoolVar(&stripPaths, "strip", false, "Strip directory paths to and just use base filenames")
}
