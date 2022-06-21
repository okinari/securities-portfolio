package util

import (
	"fmt"
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

type Stock struct {
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

func GetSecuritiesCompanyName(securitiesCompany SecuritiesCompany) string {
	switch securitiesCompany {
	case SbiSecurities:
		return "SBI証券"
	case SbiNeomobile:
		return "SBIネオモバ証券"
	case RakutenSecurities:
		return "楽天証券"
	}
	return ""
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
	if strings.Contains(str, "SBIネオモバ") {
		return SpecificAccount
	}

	return NoneAccount
}

func GetSecuritiesAccountName(securitiesAccount SecuritiesAccount) string {
	switch securitiesAccount {
	case NisaAccount:
		return "NISA口座"
	case SpecificAccount:
		return "特定口座"
	case GeneralAccount:
		return "一般口座"
	}
	return ""
}

func GetCountryName(country Country) string {
	switch country {
	case Japan:
		return "日本"
	case Usa:
		return "米国"
	}
	return ""
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
	f, err = strconv.ParseFloat(fmt.Sprintf("%.2f", f), 0)
	if err != nil {
		return 0.0, err
	}
	return f, nil
}

func DiffStocks(stocksMain, stocksSub []Stock) []Stock {
	var diffStocks []Stock
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

func PrintStock(stock Stock) {
	fmt.Printf("%v, %v, %v, %v, %v, %v\n",
		GetSecuritiesCompanyName(stock.SecuritiesCompany),
		GetCountryName(stock.StockCountry),
		GetSecuritiesAccountName(stock.SecuritiesAccount),
		stock.SecuritiesCode,
		stock.NumberOfOwnedStock,
		stock.AveragePurchasePrice,
	)
}
