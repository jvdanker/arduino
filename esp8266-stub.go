/*
 * This application can be used to experiment and test various serial port options
 */

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jacobsa/go-serial/serial"
	"bufio"
	log "github.com/Sirupsen/logrus"
	"io"
	"strings"
)

func usage() {
	fmt.Println("go-serial-test usage:")
	flag.PrintDefaults()
	os.Exit(-1)
}

func main() {
	fmt.Println("Go serial test")
	port := "/dev/tty.usbmodem1411"
	baud := flag.Uint("baud", 115200, "Baud rate")
	even := flag.Bool("even", false, "enable even parity")
	odd := flag.Bool("odd", false, "enable odd parity")
	rs485 := flag.Bool("rs485", false, "enable RS485 RTS for direction control")
	rs485HighDuringSend := flag.Bool("rs485_high_during_send", false, "RTS signal should be high during send")
	rs485HighAfterSend := flag.Bool("rs485_high_after_send", false, "RTS signal should be high after send")
	stopbits := flag.Uint("stopbits", 1, "Stop bits")
	databits := flag.Uint("databits", 8, "Data bits")
	chartimeout := flag.Uint("chartimeout", 100, "Inter Character timeout (ms)")
	minread := flag.Uint("minread", 0, "Minimum read count")

	flag.Parse()

	if *even && *odd {
		fmt.Println("can't specify both even and odd parity")
		usage()
	}

	parity := serial.PARITY_NONE

	if *even {
		parity = serial.PARITY_EVEN
	} else if *odd {
		parity = serial.PARITY_ODD
	}

	options := serial.OpenOptions{
		PortName:               port,
		BaudRate:               *baud,
		DataBits:               *databits,
		StopBits:               *stopbits,
		MinimumReadSize:        *minread,
		InterCharacterTimeout:  *chartimeout,
		ParityMode:             parity,
		Rs485Enable:            *rs485,
		Rs485RtsHighDuringSend: *rs485HighDuringSend,
		Rs485RtsHighAfterSend:  *rs485HighAfterSend,
	}

	f, err := serial.Open(options)

	if err != nil {
		fmt.Println("Error opening serial port: ", err)
		os.Exit(-1)
	} else {
		defer f.Close()
	}


	for {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			s := scanner.Text()
			log.Info("PC: <<< " + s)

			if (s == "AT+RST") {
				send(f, "ready")
			}
			if (s == "AT+CWMODE=1") {
				send(f, "\n\n\n")
			}
			if (s == "AT+CIPMUX=1") {
				send(f, "\n\n\n")
			}
			if (s == "AT+CWJAP=\"xxxx\",\"xxxxxxxx\"") {
				send(f, "OK")
			}
			if (s == "AT+CIPSTART=0,\"TCP\",\"23.203.214.89\",80") {
				send(f, "OK")
			}
			if (s == "AT+CIPSTATUS") {
				send(f, "OK")
			}
			if (strings.HasPrefix(s, "AT+CIPSEND=0,")) {
				send(f, ">")
			}
			if (strings.HasPrefix(s, "POST ")) {
				send(f, "OK")
			}
		}

		// check for errors
		if err = scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
}

func send(f io.ReadWriteCloser, message string) {
	count, err := f.Write([]byte(message))
	if err != nil {
		fmt.Println("Error writing to serial port: ", err)
	} else {
		fmt.Printf("Wrote %v bytes\n", count)
	}
}