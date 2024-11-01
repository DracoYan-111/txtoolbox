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
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	utils "txtoolbox/cmd/utils"

	"github.com/common-nighthawk/go-figure"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// TransactionCmd represents the transaction/transaction command
var TransactionCmd = &cobra.Command{
	Use:   "trade",
	Short: "Use shell to initiate transactions on blockchain directly",
	Long:  figure.NewFigure("trade", "", true).String(),
	Run: func(cmd *cobra.Command, args []string) {
		readInConfig()
		fmt.Println("transaction/trade called")
	},
}

func init() {
	// Add command
}

type Trade struct {
	NetWork     string
	ChainId     *big.Int
	FromAddress common.Address
	Private     string
	To          *common.Address
	Amount      string
	Nonce       uint64
	GasPrice    *big.Int
	GasLimit    uint64
	Data        []byte
}

// Reading Configuration Files
func readInConfig() *Trade {
	trade := new(Trade)
	trade.NetWork = viper.GetString("netWork")
	trade.Private = viper.GetString("privateKey")
	trade.To = new(common.Address)
	*trade.To = common.HexToAddress(viper.GetString("to"))

	amount := viper.GetString("amount")
	if amount == "" {
		trade.Amount = "0"
	} else {
		_, ok := new(big.Rat).SetString(viper.GetString("amount"))
		if !ok {
			trade.Amount = "0"
		} else {
			trade.Amount = amount
		}
	}
	trade.Nonce = viper.GetUint64("nonce")
	trade.GasPrice, _ = new(big.Int).SetString(viper.GetString("gasprice"), 10)
	trade.GasLimit = viper.GetUint64("gaslimit")

	trade.Data = []byte(viper.GetString("data"))

	err := processConfig(trade)
	if err != nil {
		fmt.Println(err)
	}
	return trade
}

// Processing Configuration Files
func processConfig(trade *Trade) error {
	// Check network
	if trade.NetWork == "" {
		return errors.New("netWork is empty")
	}

	client, err := ethclient.Dial(trade.NetWork)
	if err != nil {
		return err
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return err
	}
	trade.ChainId = chainID
	fmt.Println("<-- ⛓️  Network connection successful, chainID:", chainID, "-->")

	// Check privateKey
	if trade.Private == "" {
		return errors.New("privateKey is empty")
	}

	private := strings.TrimLeft(trade.Private, "0x")
	privateKey, err := crypto.HexToECDSA(private)
	if err != nil {
		return err
	}
	privateToAddr := crypto.PubkeyToAddress(privateKey.PublicKey)
	privateToAddrColor, _ := utils.GenAddressColor(privateToAddr.String())
	trade.FromAddress = privateToAddr
	fmt.Println("<-- 🥷  Private key configuration successful:", privateToAddrColor, "-->")

	// Check to address
	to := trade.To.String()
	if to == new(common.Address).String() || len(to) != 42 {
		return errors.New("to address is empty")
	}
	toAddrColcor, _ := utils.GenAddressColor(to)
	fmt.Println("<-- 💸 To Address configuration successful:", toAddrColcor, "-->")

	// Check uints and amount
	amountUints := viper.GetString("amountUint")
	if amountUints == "" || amountUints == "wei" {
		if strings.Contains(trade.Amount, ".") {
			return errors.New("the default unit is wei, and decimals are displayed for amounts")
		}
	}
	if utils.UnitMultipliers[amountUints] == "" {
		if trade.Amount != "0" {
			fmt.Println("<-- 💵 Amount Configuration Successful:", trade.Amount, "-->")
		}
	} else {
		if trade.Amount != "0" {
			uintsMap := utils.EthNumberConverter(trade.Amount, amountUints)
			fmt.Println("╔══[ 💵 Amount Configuration Successful ]══╗")
			fmt.Printf("  %-6s: %v wei\n", "wei", uintsMap["wei"])
			fmt.Printf("  %-6s: %v %s\n", amountUints, uintsMap[amountUints], amountUints)
			fmt.Println("╚══════════════════════════════════════════╝")
		}
	}

	// Check nonce
	if trade.Nonce == 0 {
		clientNonce, err := client.PendingNonceAt(context.Background(), privateToAddr)
		if err != nil {
			return err
		}
		trade.Nonce = clientNonce
	}

	// Check data
	data := viper.GetString("data")
	if trade.Data != nil {
		if strings.HasPrefix(data, "0x") {
			// Convert the input hexadecimal string to a byte array
			decodedData, err := hex.DecodeString(data[2:])
			if err != nil {
				return errors.New("failed to decode data")
			}
			trade.Data = decodedData
		}
	} else {
		trade.Data = nil
	}

	// Check gasPrice
	if trade.GasPrice == nil {
		trade.GasPrice, err = client.SuggestGasPrice(context.Background())
		if err != nil {
			return err
		}
	}

	uintsMap := utils.EthNumberConverter(trade.GasPrice.String(), "wei")
	fmt.Println("╔═[ 💰 GasPrice configuration successful ]═╗")
	fmt.Printf("  %-6s: %v wei\n", "wei", trade.GasPrice.String())
	fmt.Printf("  %-6s: %v %s\n", "gwei", uintsMap["gwei"], "gwei")
	fmt.Println("╚══════════════════════════════════════════╝")

	// Check gasLimit
	if trade.GasLimit == 0 {
		gasLimit, err := estimateTxGas(client, trade)
		if err != nil {
			return err
		}
		trade.GasLimit = gasLimit
	}

	gasLimitString := strconv.FormatUint(trade.GasLimit, 10)
	uintsMap = utils.EthNumberConverter(gasLimitString, "gwei")
	fmt.Println("╔═[ 🏦 GasLimit configuration successful ]═╗")
	fmt.Printf("  %-6s: %v wei\n", "wei", uintsMap["wei"])
	fmt.Printf("  %-6s: %v %s\n", "gwei", uintsMap["gwei"], "gwei")
	fmt.Println("╚══════════════════════════════════════════╝")

	fmt.Println("<-- 🪤  Nonce configuration successful:", trade.Nonce, "-->")

	if len(trade.Data) > 0 {
		fmt.Println("<-- 📝 Data configuration successful:", data, "-->")
	}

	for {
		fmt.Println("Start transaction? (Y/y/N/n)")

		var next string
		fmt.Scanln(&next)

		switch next {
		case "Y", "y":
			err = initiateTx(client, trade)
			if err != nil {
				return err
			}
			return nil
		case "N", "n":
			os.Exit(0)
		default:
			continue
		}
	}
}

// Estimate the gas limit
func estimateTxGas(client *ethclient.Client, trade *Trade) (uint64, error) {
	value, _ := new(big.Int).SetString(trade.Amount, 10)

	callMsg := ethereum.CallMsg{
		From:     trade.FromAddress,
		To:       trade.To,
		GasPrice: trade.GasPrice,
		Value:    value,
		Data:     trade.Data,
	}

	gasLimit, err := client.EstimateGas(context.Background(), callMsg)
	if err != nil {
		return uint64(0), err
	}
	return gasLimit, nil
}

// Initiate a transaction
func initiateTx(client *ethclient.Client, trade *Trade) error {

	// Create the transaction
	amount, _ := new(big.Int).SetString(trade.Amount, 10)
	tx := types.NewTransaction(trade.Nonce, *trade.To, amount, trade.GasLimit, trade.GasPrice, trade.Data)

	// Sing the transaction
	private := strings.TrimLeft(trade.Private, "0x")
	privateKey, _ := crypto.HexToECDSA(private)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(trade.ChainId), privateKey)
	if err != nil {
		return errors.New("signature transaction failed")
	}

	fmt.Println("<-- 📝 Tx hash configuration successful:", signedTx.Hash().Hex(), "-->")

	for {
		fmt.Println("Send transaction? (Y/y/N/n)")

		var next string
		fmt.Scanln(&next)

		switch next {
		case "Y", "y":
			// Send the transaction
			err = client.SendTransaction(context.Background(), signedTx)
			if err != nil {
				return err
			}
			fmt.Println("<-- 🚀 Transaction sent-->")
			return nil
		case "N", "n":
			os.Exit(0)
		default:
			continue
		}
	}
}
