package util

import (
	"time"
)

type SecuritiesCompany int

const (
	NoneSecuritiesCompany SecuritiesCompany = iota
	SbiSecurities
	SbiNeomobile
)

type SecuritiesAccount int

const (
	NoneAccount SecuritiesAccount = iota
	NisaAccount
	SpecificAccount
)

type StockInfo struct {
	SecuritiesCompany    SecuritiesCompany
	SecuritiesAccount    SecuritiesAccount
	SecuritiesCode       int
	AveragePurchasePrice int
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
