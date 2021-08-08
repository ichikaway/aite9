package main

import (
	"aite9/printer"
	"bufio"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

var VERSION = "0.0.1"

func showUsage() {
	printer.ErrorPrintf("USAGE:  go run aite9.go -tcp '22,3306,8080' -mode slient < serverlist.txt \n")
	os.Exit(1)
}

func main() {

	fp := os.Stdin
	serverList := getServerList(fp)


	if len(serverList) == 0 {
		showVersion()
		showUsage()
	}

	tcpPortsString := flag.String("tcp", "22", "tcp port list with comma. ex. 80,443,111")
	mode := flag.String("mode", "", "optional: silent")
	flag.Parse()

	if *mode == "silent" {
		printer.SilentModeOn()
	}

	infoDump(*tcpPortsString)
	portList := parsePortList(*tcpPortsString)

	var errorCount int = 0
	for _, server := range serverList {
		if (lookupCheck(server) == false) {
			printer.Printf("lookup failed: " + server + "\n")
			continue
		}
		for _, port := range portList {
			printer.Printf("start scan: " + server + ":" + port + "\n")
			result, err := scanOpenPort(server, port)
			if (result) {
				errorCount++
				fmt.Println(err)
			}
		}
	}

	fmt.Printf("=== Result: error count %d \n", errorCount)
	if errorCount > 0 {
		os.Exit(1)
	}
}

func lookupCheck(server string) bool {
	_,err := net.LookupHost(server)
	if err == nil {
		return true
	}
	return false
}

func scanOpenPort(server string, port string) (bool, error) {
	address := server + ":" + port
	conn, err := net.DialTimeout("tcp", address, time.Duration(1) * time.Second)
	if err == nil {
		conn.Close()
		//port open. It's danger
		return true, errors.New(" open port error: " + address)
	}
	return false, nil
}

func parsePortList(portString string) []string {
	portStringList := strings.Split(portString, ",")
	return portStringList
}

func getServerList(fp *os.File) []string {
	serverList := []string{}

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Index(line, "#") >= 0 {
			list := strings.Split(line, "#")
			server := strings.Trim(list[0], " ")
			if server == "" {
				continue
			}
			line = server
		}
		serverList = append(serverList, line)
	}
	return serverList
}

func infoDump(tcpPorts string) {
	showVersion()
	printer.Printf("- TCP port list: %s\n\n", tcpPorts)
	printer.Printf("")
}

func showVersion() {
	printer.Printf("=== aite9 Version: %s === \n", VERSION)
}
