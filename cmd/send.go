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
	"fmt"
	"io"
	"net"
	"os"

	"github.com/vikulin/go-udt/udt"
	"github.com/oliverisaac/goli"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.SetArgs([]string{"/dev/stdin"})
			return nil
		}
		for i, a := range args {
			if a == "-" {
				args[i] = "/dev/stdin"
			}
		}
		cmd.SetArgs(args)
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		goli.InitLogrus(logrus.DebugLevel)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		remoteDialAddr := fmt.Sprintf("%s:%d", destAddress, destPort)
		remoteAddr, err := net.ResolveUDPAddr("udp4", remoteDialAddr)
		if err != nil {
			return errors.Wrap(err, "Resolving a local UDP port to send from")
		}

		conn, err := udt.DialUDT("udp4", "127.0.0.1:9898", remoteAddr, false)
		if err != nil {
			return errors.Wrapf(err, "Dialing from :0 to %s", &remoteDialAddr)
		}
		totalWritten := 0
		for {
			bytesWritten, err := io.CopyN(conn, os.Stdin, MESSAGE_SIZE)
			totalWritten += int(bytesWritten)
			logrus.Debugf("Wrote %d bytes (total: %d)", bytesWritten, totalWritten)
			if err != nil {
				if err.Error() == "EOF" {
					logrus.Infof("Done writing file")
					break
				}
				return errors.Wrapf(err, "Failed to write to destination")
			}
		}
		logrus.Debug("Closing conn")
		conn.Close()
		logrus.Debug("Done closing conn")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)

	sendCmd.Flags().StringVarP(&destAddress, "address", "a", "", "Address to send to")
	sendCmd.MarkFlagRequired("address")

	sendCmd.Flags().IntVarP(&destPort, "port", "p", 9876, "Port to send to")
	sendCmd.Flags().BoolVar(&absolutePaths, "absolute", false, "Strip directory paths to and just use base filenames")
}
