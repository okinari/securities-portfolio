package util

import (
	"strconv"
	"strings"
	"time"
)

type SecuritiesCompany int

const (
	NoneSecuritiesCompany SecuritiesCompany = iota
	SbiSecurities
	SbiNeomobile
	RakutenSecurities
)

type Country int

const (
	NoneCountry Country = iota
	Japan
	America
)

type SecuritiesAccount int

const (
	NoneAccount SecuritiesAccount = iota
	NisaAccount
	SpecificAccount
)

type StockInfo struct {
	SecuritiesCompany    SecuritiesCompany
	StockCountry         Country
	SecuritiesAccount    SecuritiesAccount
	SecuritiesCode       string
	AveragePurchasePrice float64
	NumberOfOwnedStock   int
}

func GetSecuritiesCompany(str string) SecuritiesCompany {
	if str == "SBI証券" {
		return SbiSecurities
	}
	if str == "SBIネオモバ" {
		return SbiNeomobile
	}

	return NoneSecuritiesCompany
}

func WaitTime() {
	time.Sleep(3 * time.Second)
}

func ToIntByRemoveString(str string) int {
	n := 0
	for _, r := range str {
		if '0' <= r && r <= '9' {
			n = n*10 + int(r-'0')
		}
	}
	return n
}

func ToFloatByRemoveString(str string) (float64, error) {
	strNumber := ""
	slice := strings.Split(str, "")
	for i, r := range str {
		if ('0' <= r && r <= '9') || r == '.' {
			strNumber += slice[i]
		}
	}
	f, err := strconv.ParseFloat(strNumber, 0)
	if err != nil {
		return 0.0, err
	}
	return f, nil
}
