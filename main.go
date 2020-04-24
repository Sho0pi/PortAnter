package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	CONNECTION_TIMEOUT = 2
)

const (
	colorWhite  = "\033[37m"
	colorReset  = "\033[0m"

	colorPurple = "\033[35m"
	colorYellow = "\033[33m"
)

func printStartMessage() {
	fmt.Printf("%s______          _    ___        _            		%s \\       /\n", colorWhite, colorPurple)
	fmt.Printf("%s| ___ \\        | |  / _ \\      | |           		%s  \\     / \n", colorWhite, colorPurple)
	fmt.Printf("%s| |_/ /__  _ __| |_/ /_\\ \\_ __ | |_ ___ _ __ 		%s   \\.-./\n", colorWhite, colorPurple)
	fmt.Printf("%s|  __/ _ \\| '__| __|  _  | '_ \\| __/ _ \\ '__|		%s  (o\\^/o)  _   _   _     __\n", colorWhite, colorPurple)
	fmt.Printf("%s| | | (_) | |  | |_| | | | | | | ||  __/ |   		%s   ./ \\.\\ ( )-( )-( ) .-'  '-.\n", colorWhite, colorPurple)
	fmt.Printf("%s\\_|  \\___/|_|   \\__\\_| |_/_| |_|\\__\\___|_|   		%s    {-} \\(//  ||   \\\\/ (   )) '-.\n",colorWhite, colorPurple)
	fmt.Printf("							%s         //-__||__.-\\\\.       .-'\n", colorPurple)
	fmt.Printf("							%s        (/    ()     \\)'-._.-'\n", colorPurple)
	fmt.Printf("%sPlease enter a IP/URL you want to scan:				||    ||      \\\\%s\n", colorPurple, colorReset)
	fmt.Println()
}

func isIn(slice []string, value string) bool {
	for _, v := range slice {
		if value == v {
			return true
		}
	}
	return false
}

func getServerIps(url string) ([]string, error) {
	var ans []string
	ips, err := net.LookupIP(url)
	if err != nil {
		return nil, err
	}
	for _, ip := range ips {
		ipv4 := ip.To4().String()
		if !isIn(ans, ipv4) {
			ans = append(ans, ipv4)
		}
	}
	return ans, nil
}

func getUserInputToIps() ([]string, error) {
	fmt.Print("")
	var input string
	for {
		fmt.Scanf("%s", &input)
		if ips, err := getServerIps(input); err == nil {
			return ips, nil
		} else {
			if _, ok := err.(*net.DNSError); ok {
				fmt.Print("Please try again: ")
			} else {
				return nil, err
			}
		}

	}
}

func main() {
	printStartMessage()

	ips, err := getUserInputToIps()
	if err != nil {
		panic(err)
	}
	for _, ip := range ips {
		fmt.Println(ip)
	}
	var wg sync.WaitGroup

	for _, ip := range ips {
		for port := 1; port <= 65535; port++ {
			wg.Add(1)
			time.Sleep(20 * time.Millisecond)
			go checkPort(ip, port, &wg)

		}
	}
	wg.Wait()
}

func checkPort(ip string, port int, wg *sync.WaitGroup) {
	defer wg.Done()
	addr := buildAddress(ip, port)
	for {
		_, err := net.DialTimeout("tcp", net.JoinHostPort(ip, strconv.Itoa(port)), CONNECTION_TIMEOUT*time.Second)
		if err == nil {
			fmt.Println(addr)
			return
		}
		if nerr, ok := err.(*net.OpError); ok {
			if isTimeoutErr(nerr.Err) || isNoMoreSockets(nerr.Err) {
				time.Sleep(150 * time.Millisecond)
			} else {
				return
			}
		}
	}

}

func isTimeoutErr(err error) bool {
	var builder strings.Builder
	fmt.Fprintf(&builder, "%T", err)
	return strings.Contains(builder.String(), "TimeoutError")
}

func isNoMoreSockets(err error) bool {
	var builder strings.Builder
	fmt.Fprintf(&builder, "%s", err)
	return strings.Contains(builder.String(), "too many open files")
}

func buildAddress(ip string, port int) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "%s:%d", ip, port)
	return builder.String()
}
