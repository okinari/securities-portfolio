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
	Usa
)

type SecuritiesAccount int

const (
	NoneAccount SecuritiesAccount = iota
	NisaAccount
	SpecificAccount
	GeneralAccount
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
	if strings.Contains(str, "ネオモバ") {
		return SbiNeomobile
	}
	if strings.Contains(str, "SBI") {
		return SbiSecurities
	}
	if strings.Contains(str, "楽天") {
		return RakutenSecurities
	}

	return NoneSecuritiesCompany
}

func GetSecuritiesAccount(str string) SecuritiesAccount {
	if strings.Contains(str, "NISA") {
		return NisaAccount
	}
	if strings.Contains(str, "特定") {
		return SpecificAccount
	}
	if strings.Contains(str, "一般") {
		return GeneralAccount
	}

	return NoneAccount
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
	for _, s := range slice {
		if s == "." {
			strNumber += s
			continue
		}

		_, err := strconv.Atoi(s)
		if err != nil {
			continue
		}
		strNumber += s
	}
	f, err := strconv.ParseFloat(strNumber, 0)
	if err != nil {
		return 0.0, err
	}
	return f, nil
}

func DiffStocks(stocksMain, stocksSub []StockInfo) []StockInfo {
	var diffStocks []StockInfo
	for _, stockMain := range stocksMain {
		isNotExist := true
		for _, stockSub := range stocksSub {
			if stockSub == stockMain {
				isNotExist = false
				break
			}
		}
		if isNotExist {
			diffStocks = append(diffStocks, stockMain)
		}
	}
	return diffStocks
}
