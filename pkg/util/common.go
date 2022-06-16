package util

import "time"

type SecuritiesCompanyCode int

const (
	None SecuritiesCompanyCode = iota
	SbiSecurities
	SbiNeomobile
)

type StockInfo struct {
	SecuritiesCompany    SecuritiesCompanyCode
	SecuritiesCode       int
	AveragePurchasePrice int
	NumberOfOwnedStock   int
}

func GetSecuritiesCompanyCode(str string) SecuritiesCompanyCode {
	if str == "SBI証券" {
		return SbiSecurities
	}
	if str == "SBIネオモバ" {
		return SbiNeomobile
	}

	return None
}

func WaitTime() {
	time.Sleep(3 * time.Second)
}
