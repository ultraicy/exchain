package types

import (
	"fmt"
	"testing"
)

func TestValidateDenom(t *testing.T) {
	//0.000652393267892272ammswap_okt_usdt-a2b,1.275215023723700522okt
	err := ValidateDenom("0.000652393267892272ammswap_okt_usdt-a2b")
	fmt.Println(err)
	err = ValidateDenom("1.275215023723700522okt")
	fmt.Println(err)
}
