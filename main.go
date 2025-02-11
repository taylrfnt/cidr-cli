package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	/* GET CIDR RANGE FROM INPUT */
	var cidrBlock string
	var testIp string
	flag.StringVar(&cidrBlock, "C", "", "The CIDR Address Block in format <address>/<network prefix>")
	flag.StringVar(&testIp, "t", "", "The IP address to be tested within the CIDR range, in IPv4 format (xxx.xxx.xxx.xxx)")
	flag.Parse()

	// split the cidr block into address & network prefix
	addressSlice := strings.Split(cidrBlock, "/")
	networkPrefix, err := strconv.Atoi(addressSlice[len(addressSlice)-1])
	if err != nil {
		panic(err)
	}

	// let's split the address into bits
	decimalIp := addressSlice[0]
	ipSlice := strings.Split(decimalIp, ".")

	/* CONVERT ADDRESS DECIMAL TO BINARY */

	// let's make the IP []string into []int, then convert to binary
	binaryIpSlice := make([]string, len(ipSlice))
	for i, s := range ipSlice {
		// make the string an int
		intIp, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		// now make the int -> int64, then return the string representation in base
		// 2, then make it an int64 again
		int64Ip := int64(intIp)
		binaryIpString := strconv.FormatInt(int64Ip, 2)
		if len(binaryIpString) < 8 {
			binaryIpSlice[i] = strings.Repeat("0", 8-len(binaryIpString)) + binaryIpString
		} else {
			binaryIpSlice[i] = binaryIpString
		}
	}
	fmt.Printf("Binary IP representation:\n%s %s %s %s\n",
		binaryIpSlice[0],
		binaryIpSlice[1],
		binaryIpSlice[2],
		binaryIpSlice[3],
	)

	/* CONVERT TO START AND END ADDRESS */
	binaryIp := strings.Join(binaryIpSlice, "")
	addressPool := binaryIp[networkPrefix:]

	if strings.Contains(addressPool, "1") {
		// TODO: handle networks not ending in all zeroes
	}

	firstAddressBit := strings.Repeat("0", len(addressPool))
	lastAddressBit := strings.Repeat("1", len(addressPool))
	binaryFirstAddr := binaryIp[:networkPrefix] + firstAddressBit
	binaryLastAddr := binaryIp[:networkPrefix] + lastAddressBit

	/* now we need to split into bits, then convert back to decimal/base10 */
	firstAddrSlice := [4]string{binaryFirstAddr[:8], binaryFirstAddr[8:16], binaryFirstAddr[16:24], binaryFirstAddr[24:32]}
	lastAddrSlice := [4]string{binaryLastAddr[:8], binaryLastAddr[8:16], binaryLastAddr[16:24], binaryLastAddr[24:32]}

	decimalFirstAddr := make([]int64, len(firstAddrSlice))
	decimalLastAddr := make([]int64, len(lastAddrSlice))
	for i := range firstAddrSlice {
		decimalFirstAddr[i], err = strconv.ParseInt(firstAddrSlice[i], 2, 64)
		if err != nil {
			panic(err)
		}
		decimalLastAddr[i], err = strconv.ParseInt(lastAddrSlice[i], 2, 64)
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("First Address: %v.%v.%v.%v\nLast Address: %v.%v.%v.%v\n",
		decimalFirstAddr[0], decimalFirstAddr[1], decimalFirstAddr[2], decimalFirstAddr[3],
		decimalLastAddr[0], decimalLastAddr[1], decimalLastAddr[2], decimalLastAddr[3])
}
