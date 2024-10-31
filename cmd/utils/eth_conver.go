/*
Copyright © 2024 DracoYan-111 <yanlong2944@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
)

// EthConverCmd represents the utils/ethConver command
var EthConverCmd = &cobra.Command{
	Use:   "ethConver",
	Short: "Unit conversion on Ethereum",
	Long:  figure.NewFigure("EthConver", "", true).String(),
	Example: `
Enter the quantity and unit to be converted, such as one of the following units:
-n number -u wei
-n number -u kwei
-n number -u mwei
-n number -u gwei
-n number -u szabo
-n number -u finney
-n number -u ether
-n number -u kether
-n number -u mether
-n number -u gether
-n number -u tether
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("utils/ethConver called")

		return checkInput()
	},
}

var number string
var unit string

func init() {
	// Add flags
	EthConverCmd.Flags().StringVarP(&number, "number", "n", "", "number")
	EthConverCmd.Flags().StringVarP(&unit, "unit", "u", "", "unit")

	// Mark flags required
	EthConverCmd.MarkFlagRequired("number")
	EthConverCmd.MarkFlagRequired("unit")
}

// Check input
func checkInput() error {
	// Check number
	if isNumber := func(_number string) bool {
		_, ok := new(big.Rat).SetString(_number)
		return ok
	}; !isNumber(number) {
		return errors.New("Check the number entered:<" + number + ">")
	}

	// Check unit
	if isUint := func(_uints string) bool {
		_, ok := UnitMultipliers[_uints]
		return ok
	}; !isUint(unit) {
		return errors.New("Check the units entered:<" + unit + ">")
	}

	// Convert
	converResults := EthNumberConverter(number, unit)

	// Print the results in order
	for _, v := range UintsList {
		fmt.Printf("%-7s: %s\n", v, converResults[v])
	}
	return nil
}

// Convert input to eth units
func EthNumberConverter(number, unit string) map[string]string {
	converResults := make(map[string]string)

	// Convert number to big.Rat
	numberBigRat, _ := new(big.Rat).SetString(number)

	// Convert uint to big.Rat ready for calculation
	uintsBigInt, _ := new(big.Int).SetString(UnitMultipliers[unit], 10)
	uintsBigRat := new(big.Rat).SetInt(uintsBigInt)

	// Calculate input results
	numberInUintRat := new(big.Rat).Mul(numberBigRat, uintsBigRat)

	for _, value := range UintsList {
		// Get the precision of each unit
		uintNameValueBigInt, _ := new(big.Int).SetString(UnitMultipliers[value], 10)
		uintNameValueBigRat := new(big.Rat).SetInt(uintNameValueBigInt)

		// Calculate the result of each unit
		resultRat := new(big.Rat).Quo(numberInUintRat, uintNameValueBigRat)
		// Format the result
		result := resultRat.FloatString(30)

		//	Check if the result contains "."
		if strings.Contains(result, ".") {
			result = strings.TrimRight(result, "0")
			result = strings.TrimRight(result, ".")

			converResults[value] = result
		}
	}
	return converResults
}
