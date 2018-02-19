// example.go
package main

import (
	"fmt"
	"github.com/thanhpk/baseconv"
)

func main() {
	valHex := "70B1D707EAC2EDF4C6389F440C7294B51FFF57BB"
	valDec, _ := baseconv.DecodeHex(valHex)
	val62, _ := baseconv.Convert(valHex, baseconv.DigitsHex, baseconv.Digits62)
	val36, _ := baseconv.Encode36(val62, baseconv.Digits62)

	fmt.Println("dec string: " + valDec)
	fmt.Printf("62 string: " + val62)
	fmt.Println("36 string: " + val36)

	conVal36, _ := baseconv.Decode36(val36, baseconv.DigitsDec)
	fmt.Printf("dec and 36 values same: %t\n", valDec == conVal36)

	val62, _ = baseconv.Convert("7(n42-&DG$MT", baseconv.ASCII, baseconv.Digits62)
	str, _ := baseconv.Convert(val62, baseconv.Digits62, baseconv.ASCII)
	fmt.Println("62 string: " + val62 + " str: " + str)
}
