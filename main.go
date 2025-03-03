package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/* UTILITY FUNCTIONS */
func ParseCidr(cidrBlock string) (string, int, error) {
	splitAddr := strings.Split(cidrBlock, "/")

	prefixBits, err := strconv.Atoi(splitAddr[len(splitAddr)-1])
	if err != nil {
		fmt.Println("Error parsing CIDR block:", err)
		return "", 32, err
	} else {
		return splitAddr[0], prefixBits, nil
	}
}

func ConvertAddr(address string, targetFormat string) (string, error) {
	bitwiseAddress := strings.Split(address, ".")
	if targetFormat == "binary" {
		// placeholder for storing converted vals
		binaryAddress := make([]string, len(bitwiseAddress))
		// perform operations on each address bit to convert to binary
		for i, bit := range bitwiseAddress {
			// make the string an int
			decimalBit, err := strconv.Atoi(bit)
			if err != nil {
				return "", err
			} else {
				// convert to base 2
				binaryBit := strconv.FormatInt(int64(decimalBit), 2)
				// pad leading zeroes where the binary representation is less than
				// 8 bits.
				if len(binaryBit) < 8 {
					binaryAddress[i] = strings.Repeat("0", 8-len(binaryBit)) + binaryBit
				} else {
					binaryAddress[i] = binaryBit
				}
			}
		}
		return strings.Join(binaryAddress, "."), nil
	} else if targetFormat == "decimal" {
		// placeholder for storing values
		decimalAddress := make([]string, len(bitwiseAddress))
		for i, bit := range bitwiseAddress {
			decimalBit, err := strconv.ParseInt(bit, 10, 0)
			if err != nil {
				return "", err
			} else {
				decimalAddress[i] = strconv.Itoa(int(decimalBit))
			}
		}
		return strings.Join(decimalAddress, "."), nil
	} else {
		message := fmt.Errorf("Invalid target format supplied\nSupplied: %s\nExpected: [decimal, binary]\n", targetFormat)
		return "", message
	}
}

func NewConvertCommand() *ConvertCommand {
	// make a convert command object
	cc := &ConvertCommand{
		fs: flag.NewFlagSet("convert", flag.ContinueOnError),
	}
	// define our flags
	cc.fs.StringVar(&cc.cidr, "c", "", "The IPv4 CIDR range to convert in standard format (<address>/<network prefix>).")
	cc.fs.BoolVar(&cc.binary, "b", false, "Display conversion information in binary format (default = decimal)")

	return cc
}

type ConvertCommand struct {
	fs     *flag.FlagSet
	cidr   string
	binary bool
}

// type MemberCommand struct {
// 	fs      *flag.FlagSet
// 	address string
// 	cidr    string
// }

func (c *ConvertCommand) Name() string {
	return c.fs.Name()
}

// func (m *MemberCommand) Name() string {
// 	return m.fs.Name()
// }

func (c *ConvertCommand) Init(args []string) error {
	return c.fs.Parse(args)
}

// func (m *MemberCommand) Init(args []string) error {
// 	return m.fs.Parse(args)
// }

func (c *ConvertCommand) Run() error {
	address, networkPrefix, err := ParseCidr(c.cidr)
	if err != nil {
		panic(err)
	}

	/* CONVERT ADDRESS DECIMAL TO BINARY */
	binaryAddress, err := ConvertAddr(address, "binary")

	/* CONVERT TO START AND END ADDRESS

	Now that we have a binary representation of our address in 32 bits, we can convert this to a
	start and end IP address to define the range this CIDR block corresponds to.

	We'll (again) do this with stdlib modules and strings, because I am a masochist and wish this
	(albeit mild) complexity upon myself.

	We need to:
	1. Use the network prefix to determine the bits that comprise the address addressPool
	2. Retain the bits up to the network prefix/subnet mask
	3. Replace the address pool with all zeroes and all ones, then tack those back onto the fixed
	bits to compose a first & last address in the network range.

	There's probably room to debate on if we should display a range starting with a `.0` or `.1`
	(and by that same line of thought, a range ending in `.255` or `.254`) due to multicast &
	broadcast, but we'll leave that to the user to understand and just display the first (`.0`) and
	last (`.255`) addresses in our range conversion.
	*/

	// take the bits AFTER the network prefix for modification
	binaryBits := strings.Split(binaryAddress, ".")
	binaryAddress = strings.Join(binaryBits, "")
	addressPool := binaryAddress[networkPrefix:]

	binaryFirstAddr := binaryAddress[:networkPrefix] + strings.Repeat("0", len(addressPool))
	binaryLastAddr := binaryAddress[:networkPrefix] + strings.Repeat("1", len(addressPool))

	/* now we need to split into bits, then convert back to decimal/base10 */
	binaryFirstAddrSlice := [4]string{binaryFirstAddr[:8], binaryFirstAddr[8:16], binaryFirstAddr[16:24], binaryFirstAddr[24:32]}
	binaryLastAddrSlice := [4]string{binaryLastAddr[:8], binaryLastAddr[8:16], binaryLastAddr[16:24], binaryLastAddr[24:32]}

	decimalFirstAddr := make([]int64, len(binaryFirstAddrSlice))
	decimalLastAddr := make([]int64, len(binaryLastAddrSlice))
	// for IPv4, this will always be 4
	for i := range decimalFirstAddr {
		decimalFirstAddr[i], err = strconv.ParseInt(binaryFirstAddrSlice[i], 2, 64)
		if err != nil {
			panic(err)
		}
		decimalLastAddr[i], err = strconv.ParseInt(binaryLastAddrSlice[i], 2, 64)
		if err != nil {
			panic(err)
		}
	}

	/* PRINTING OUTPUT TO CONSOLE
	Now that we have performed all our conversions and have stored results in vars, we can print the
	value to the console for the user based on their flags.

	We will always print the values in decimal (base 10) representation.  If the user supplies a
	flag to print binary (base 2) representation, this is displayed in addition to the decimal
	representations.

	*/

	// print first & last address
	if c.binary {
		fmt.Println("First Address                                  Last Address")
		fmt.Printf("%v.%v.%v.%v                                    %v.%v.%v.%v\n",
			decimalFirstAddr[0], decimalFirstAddr[1], decimalFirstAddr[2], decimalFirstAddr[3],
			decimalLastAddr[0], decimalLastAddr[1], decimalLastAddr[2], decimalLastAddr[3],
		)
		fmt.Printf("%v.%v.%v.%v            %v.%v.%v.%v\n",
			binaryFirstAddrSlice[0], binaryFirstAddrSlice[1], binaryFirstAddrSlice[2], binaryFirstAddrSlice[3],
			binaryLastAddrSlice[0], binaryLastAddrSlice[1], binaryLastAddrSlice[2], binaryLastAddrSlice[3],
		)
	} else {
		fmt.Println("First Address     Last Address")
		fmt.Printf("%v.%v.%v.%v       %v.%v.%v.%v\n",
			decimalFirstAddr[0], decimalFirstAddr[1], decimalFirstAddr[2], decimalFirstAddr[3],
			decimalLastAddr[0], decimalLastAddr[1], decimalLastAddr[2], decimalLastAddr[3])
	}
	return nil
}

type Runner interface {
	Init([]string) error
	Run() error
	Name() string
}

func root(args []string) error {
	if len(args) < 1 {
		return errors.New("You must pass a sub-command.")
	}

	cmds := []Runner{
		NewConvertCommand(),
	}

	subcommand := os.Args[1]

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			cmd.Init(os.Args[2:])
			return cmd.Run()
		}
	}

	return fmt.Errorf("Unknown subcommand: %s\nExpected: [ convert ]", subcommand)
}

func main() {
	if err := root(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
