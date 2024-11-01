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
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
)

// Unit precision
var UnitMultipliers = map[string]string{
	"wei":    "1",
	"kwei":   "1000",
	"mwei":   "1000000",
	"gwei":   "1000000000",
	"szabo":  "1000000000000",
	"finney": "1000000000000000",
	"ether":  "1000000000000000000",
	"kether": "1000000000000000000000",
	"mether": "1000000000000000000000000",
	"gether": "1000000000000000000000000000",
	"tether": "1000000000000000000000000000000",
}

var UintsList = []string{
	"wei",
	"kwei",
	"mwei",
	"gwei",
	"szabo",
	"finney",
	"ether",
	"kether",
	"mether",
	"gether",
	"tether",
}

// UtilsCmd represents the utils/utils command
var UtilsCmd = &cobra.Command{
	Use:   "utils",
	Short: "Different on-chain tools are available here",
	Long:  figure.NewFigure("Utils", "", true).String(),
	Example: `
ethConver -n number -u unit:Convert input to eth units
checkAddrsss -h:Different functions for addresses
`,
}

func init() {
	//Add command
	UtilsCmd.AddCommand(EthConverCmd)
	UtilsCmd.AddCommand(CheckAddressCmd)
}
