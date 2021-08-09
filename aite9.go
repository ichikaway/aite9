package main

import (
	"aite9/notification"
	"aite9/printer"
	"bufio"
	"errors"
	"flag"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var VERSION = "0.0.1"

func showUsage() {
	printer.ErrorPrintf("USAGE:  go run aite9.go -tcp '22,3306,8080' -mode silent < serverlist.txt \n")
	os.Exit(1)
}

func main() {

	fp := os.Stdin
	serverList := getServerList(fp)


	if len(serverList) == 0 {
		showVersion()
		showUsage()
	}

	tcpPortsString := flag.String("tcp", "", "tcp port list with comma. ex. 80,443,111")
	udpPortsString := flag.String("udp", "", "udp port list with comma. ex. 80,443,111")
	mode := flag.String("mode", "", "optional: silent")
	flag.Parse()

	if *mode == "silent" {
		printer.SilentModeOn()
	}

	infoDump(*tcpPortsString, *udpPortsString)
	tcpPortList := parsePortList(*tcpPortsString)
	udpPortList := parsePortList(*udpPortsString)

	var errorCount int = 0
	var errorMessageList = []string{}
	for _, server := range serverList {
		if (lookupCheck(server) == false) {
			printer.Printf("lookup failed(skip scan): " + server + "\n")
			continue
		}
		for _, port := range tcpPortList {
			if port == "" {
				continue
			}
			printer.Printf("start TCP scan: " + server + ":" + port + "\n")
			result, err := scanOpenTcpPort(server, port)
			if (result) {
				errorCount++
				errorMessageList = append(
					errorMessageList,
					err.Error(),
				)
			}
		}
		for _, port := range udpPortList {
			if port == "" {
				continue
			}
			printer.Printf("start UDP scan: " + server + ":" + port + "\n")
			result, err := scanOpenUdpPort(server, port)
			if (result) {
				errorCount++
				errorMessageList = append(
					errorMessageList,
					err.Error(),
				)
			}
		}
	}

	if errorCount > 0 {
		errorText := strings.Join(errorMessageList, "\n")
		printer.Printf("\n=== Result =====================")
		printer.Printf("\n%s\n",errorText)
		printer.Printf(" Result: error count %d \n", errorCount)
		printer.Printf("============================\n\n")

		notification.PostSlack("aite9 (Ver. "+VERSION+"), error count:" + strconv.Itoa(errorCount), errorText)
		os.Exit(1)
	} else {
		printer.Printf("\n=== Result: All OK. no error \n\n")
	}
}

func lookupCheck(server string) bool {
	_,err := net.LookupHost(server)
	if err == nil {
		return true
	}
	return false
}

func scanOpenTcpPort(server string, port string) (bool, error) {
	address := server + ":" + port
	conn, err := net.DialTimeout("tcp", address, time.Duration(1) * time.Second)
	if err == nil {
		conn.Close()
		//port open. It's danger
		printer.Printf(" open TCP port error: " + address + "\n")
		return true, errors.New(" open TCP port error: " + address + "\n")
	}
	return false, nil
}
func scanOpenUdpPort(server string, port string) (bool, error) {
	address := server + ":" + port
	conn, _ := net.DialTimeout("udp", address, time.Duration(1) * time.Second)

	writeCount := 0
	for i := 0; i < 3; i++ {
		buf := []byte("12345678912345")
		a, _ := conn.Write(buf)
		printer.Printf("\nwrite: " + strconv.Itoa(a))
		//buffer := make([]byte, 1)
		//length, err := conn.Read(buffer)
		//if length > 0 || err == nil {
		if a > 0 {
			writeCount++
		}
		conn.Close()
	}

	if writeCount > 0 {
		//port open. It's danger
		printer.Printf(" open UDP port error: " + address + "\n")
		return true, errors.New(" open UDP port error: " + address + "\n")
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

func infoDump(tcpPorts string, udpPorts string) {
	showVersion()
	printer.Printf("- TCP port list: %s\n", tcpPorts)
	if udpPorts != "" {
		printer.Printf("- UDP port list: %s\n\n", udpPorts)
	}
	printer.Printf("")
}

func showVersion() {
	printer.Printf("=== aite9 Version: %s === \n", VERSION)
}
