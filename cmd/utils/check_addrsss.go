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
	"hash/fnv"

	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
)

// CheckAddrsssCmd represents the utils/checkAddrsss command
var CheckAddrsssCmd = &cobra.Command{
	Use:   "checkAddrsss",
	Short: "Add unique colors to addresses and check for differences",
	Long:  figure.NewFigure("checkAddrsss", "", true).String(),
	Example: `	
utils color -a:Add a unique color to the address`,
}

// Add a unique color to the address
var addressColorCmd = &cobra.Command{
	Use:   "color",
	Short: "Add a unique color to the address",
	Example: `
utils color -a:Add a unique color to the address`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("utils/color called")
		color, err := genAddressColor(addressLeft)
		if err != nil {
			return err
		}
		fmt.Println(color)

		return nil
	},
}

// Compare two colors and color different characters
var colorDiffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Compare two colors and color different characters",
	Example: `
utils diff -l -r:Compare two colors and color different characters`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("utils/diff called")
		err := addrssCheckDiffCmd(addressLeft, addressRight)
		return err
	},
}

var addressLeft string
var addressRight string

func init() {
	// Add command
	CheckAddrsssCmd.AddCommand(addressColorCmd)
	CheckAddrsssCmd.AddCommand(colorDiffCmd)

	// Add flags
	addressColorCmd.Flags().StringVarP(&addressLeft, "address", "a", "", "address")
	addressColorCmd.MarkFlagRequired("address")

	colorDiffCmd.Flags().StringVarP(&addressLeft, "left", "l", "", "left address")
	colorDiffCmd.Flags().StringVarP(&addressRight, "right", "r", "", "right address")
	colorDiffCmd.MarkFlagRequired("left")
	colorDiffCmd.MarkFlagRequired("right")
}

// Use the FNV hash function to generate colors
func charToColor(char byte) string {
	// Make the colors richer with New64a
	h := fnv.New64a()
	h.Write([]byte{char})
	hash := h.Sum64()

	// Convert the hash to a color
	r := (hash & 0xFF0000) >> 16
	g := (hash & 0x00FF00) >> 8
	b := hash & 0x0000FF

	// ANSI 256 color format "\033[38;2;R;G;Bm"
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
}

// Generates a unique color for the input address
func genAddressColor(address string) (string, error) {
	// Check if the address is valid
	if len(address) != 42 {
		return "", errors.New("please enter a valid address")
	}

	var result string
	for i := 0; i < len(address); i++ {
		color := charToColor(address[i])
		result += fmt.Sprintf("%s%c\033[0m", color, address[i])
	}

	return result, nil
}

// Compare two colors and color different characters
func addrssCheckDiffCmd(addressL, addressR string) error {
	if len(addressL) != 42 || len(addressR) != 42 {
		return errors.New("please enter a valid address")
	}

	var line1, line2 string
	var difference bool
	for i := 0; i < len(addressL); i++ {
		char1 := addressL[i]
		char2 := addressR[i]
		if char1 != char2 {
			color1 := charToColor(char1)
			color2 := charToColor(char2)

			line1 += fmt.Sprintf("%s%c\033[0m", color1, char1)
			line2 += fmt.Sprintf("%s%c\033[0m", color2, char2)
			difference = true

		} else {
			line1 += fmt.Sprintf("%c", char1)
			line2 += fmt.Sprintf("%c", char2)
		}
	}
	fmt.Println("Left address -> ", line1)
	fmt.Println("Right address -> ", line2)
	fmt.Println("Difference -> ", difference)
	return nil
}
