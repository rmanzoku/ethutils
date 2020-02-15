package types

import (
	"math/big"

	"github.com/shopspring/decimal"
)

var (
	baseEther, _ = decimal.NewFromString("1000000000000000000")
	baseGwei, _  = decimal.NewFromString("1000000000")
)

type Wei big.Int

func NewWei(wei *big.Int) *Wei {
	w := Wei(*wei)
	return &w
}

func weiFromDecimal(wei decimal.Decimal) *Wei {
	r, _ := new(big.Int).SetString(wei.String(), 10)
	w := Wei(*r)
	return &w
}

func NewWeiFromString(wei string) *Wei {
	i, _ := new(big.Int).SetString(wei, 10)
	return NewWei(i)
}

func NewWeiFromEther(ether float64) *Wei {
	d := decimal.NewFromFloat(ether)
	return weiFromDecimal(d.Mul(baseEther))
}

func NewWeiFromGwei(gwei float64) *Wei {
	d := decimal.NewFromFloat(gwei)
	return weiFromDecimal(d.Mul(baseGwei))
}

func (w Wei) String() string {
	i := big.Int(w)
	return i.String()
}

func (w Wei) Decimal() decimal.Decimal {
	ret, _ := decimal.NewFromString(w.String())
	return ret
}

func (w Wei) ToGwei() float64 {
	ret, _ := w.Decimal().Div(baseGwei).Float64()
	return ret
}

func (w Wei) ToEther() float64 {
	ret, _ := w.Decimal().Div(baseEther).Float64()
	return ret
}
