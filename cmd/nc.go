/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// ncCmd represents the nc command
var ncCmd = &cobra.Command{
	Use:   "nc",
	Short: "The mytool net cat, nc, remote shell sub-tool",
	Long: `The mytool net cat, nc, sub-tool. Use ctrl-D to send EOF to end client send`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Lookup("proto").Changed {
			fmt.Println("TCP only today, sorry.")
			os.Exit(-1)
			return
		}
		if cmd.Flags().Lookup("listen").Changed {
			fmt.Println("Server mode...")
			tcpServer(host, port)
		} else {
			fmt.Println("Client mode")
			tcpClient(host, port)
		}

	},
}

var host, port, proto string
//var port uint16
var listen bool

func init() {
	rootCmd.AddCommand(ncCmd)
	ncCmd.Flags().StringVar(&host, "host", "localhost", "--host=localhost")
	ncCmd.Flags().BoolVar(&listen, "listen", false, "--listen=false")
	ncCmd.Flags().StringVar(&proto, "proto", "tcp", "--proto=tcp")
	ncCmd.Flags().StringVar(&port, "port", "8021", "--port=8021")
}

var serverConn *net.TCPConn
var connected bool

func tcpClient(host string, port string){
	_,err := net.ResolveTCPAddr("tcp", host+":"+port)
	if err!=nil {
		fmt.Println("Unable to Resolve Address", host, err)
		os.Exit(-1)
		return
	}

	ipconn, err := net.DialTimeout("tcp", host+":"+port, time.Second*5)

	if err!=nil {
		fmt.Println("Timed out dialing", host, err)
		os.Exit(-1)
		return
	}
	defer ipconn.Close()
	conn := ipconn.(*net.TCPConn)
	sendToServer(conn)
}

func sendToServer(conn *net.TCPConn){
	src := os.Stdin
	buff := make([]byte, 2048)
	for {
		var nBytes int
		var err error
		nBytes, err = src.Read(buff)
		if err!=nil {
			if err != io.EOF {
				fmt.Println("Error reading buffer", err)
			}
			os.Exit(-1)
			return
		}
		_, err = conn.Write(buff[0:nBytes])
		if err!=nil {
			fmt.Println("Error writing to Server", err)
		}
	}
}

func tcpServer(host string, port string){
	addr,err := net.ResolveTCPAddr("tcp",host+":"+port)
	if err!=nil {
		fmt.Println("Unable to Resolve Address", err)
		os.Exit(-1)
		return
	}
	lis,err := net.ListenTCP("tcp",addr)
	if err!=nil {
		fmt.Println("Unable able to Listen on TCP port", err)
		os.Exit(-1)
		return
	}
	fmt.Println("Listening on:", host, "port", port)
	for{
		conn, err := lis.AcceptTCP()
		if err!=nil {
			fmt.Println("Unable to Accept TCP", err)
			os.Exit(-1)
			return
		}
		fmt.Println("Connection from" + conn.RemoteAddr().String())
		connected = true
		serverConn = conn
		for {
			buff := make([]byte, 2048)
			buffLen, err := conn.Read(buff)
			if err!=nil {
				fmt.Println("Error reading from" + conn.RemoteAddr().String())
				connected = false
				conn.Close()
				break
			}
			if buff[0]==3 {
				fmt.Println("Error in data from" + conn.RemoteAddr().String())
				connected = false
				conn.Close()
				break
			}
			buff = buff[:buffLen]
			fmt.Print(string(buff))
		}
	}
}