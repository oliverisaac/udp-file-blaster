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
	"context"
	"net"
	"os"

	"github.com/vikulin/go-udt/udt"
	"github.com/oliverisaac/goli"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var listenAddress string
var baseDirectory string
var stripPaths bool

// receiveCmd represents the receive command
var receiveCmd = &cobra.Command{
	Use:   "receive",
	Short: "Start a client to receive files from a sender",
	Long:  `Starting a client will tell you where to send the files to`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		goli.InitLogrus(logrus.DebugLevel)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		logrus.Debugf("Opening listener on %s", listenAddress)
		udtServer := udt.DefaultConfig()
		udtServer.CanAcceptStream = false
		listener, err := udtServer.Listen(context.Background(), "udp4", listenAddress)
		if err != nil {
			return errors.Wrapf(err, "Creating UDT listener on %s", listenAddress)
		}

		defer listener.Close()

		for { // TODO: Parallelize this
			logrus.Debugf("Waiting to accept...")
			conn, err := listener.Accept() // This is a blocking call
			logrus.Debugf("Got connection")
			if err != nil {
				return errors.Wrapf(err, "Accepting a connection from the listener")
			}
			go func() {
				logrus.Debugf("Handling connection")
				err = handleReceiveConnection(conn)
				logrus.Debugf("Done handling connection")
				if err != nil {
					logrus.Error(errors.Wrapf(err, "Handling receive connection"))
				}
			}()
		}
	},
}

func handleReceiveConnection(conn net.Conn) error {
	logrus.Infof("Starting read of bytes")
	bytesRead := 0
	var numBytes int
	var err error
	for {
		logrus.Debugf("Reading from connection")
		buffer := make([]byte, MESSAGE_SIZE)
		numBytes, err = conn.Read(buffer)
		bytesRead = bytesRead + numBytes
		logrus.Debugf("Received %d bytes (total: %d)", numBytes, bytesRead)
		if err != nil {
			return errors.Wrapf(err, "After reading %d bytes", bytesRead)
		}

		logrus.Debugf("Writing to stdout")
		os.Stdout.Write(buffer[:numBytes])
	}
	logrus.Infof("Wrote all of received data: %d bytes", bytesRead)
	return nil
}

func init() {
	rootCmd.AddCommand(receiveCmd)

	receiveCmd.Flags().StringVarP(&listenAddress, "address", "a", "127.0.0.1:9876", "Address to listen on")
	receiveCmd.Flags().StringVarP(&baseDirectory, "dir", "d", "/", "Base directory to place the received files. Files will be placed relative to this directory, so if you receive file '/var/logs/messages' and specify --dir=/example then the file will land in /example/var/logs/messages")
	receiveCmd.Flags().BoolVar(&stripPaths, "strip", false, "Strip directory paths to and just use base filenames")
}
