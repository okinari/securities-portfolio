package util

import (
	"fmt"
	"regexp"
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
	SecuritiesCompany       SecuritiesCompany // 証券会社
	StockCountry            Country           // 国
	SecuritiesAccount       SecuritiesAccount // 証券口座
	SecuritiesCode          string            // 証券コード
	CompanyName             string            // 会社名
	Industry                string            // 業種
	NumberOfOwnedStock      float64           // 保有株式数
	AveragePurchasePriceOne float64           // 平均取得単価(1株)
	AveragePurchasePriceAll float64           // 平均取得単価(合計)
	ValuationOne            float64           // 評価額(1株)
	ValuationAll            float64           // 評価額(合計)
	ProfitAndLossOne        float64           // 損益(1株)
	ProfitAndLossAll        float64           // 損益(合計)
	ProfitAndLossRatio      float64           // 損益(割合)
	DividendOne             float64           // 配当金(1株)
	DividendAll             float64           // 配当金(合計)
	DividendRatio           float64           // 配当利回り
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
		if ('0' <= r && r <= '9') || (r == '-') {
			n = n*10 + int(r-'0')
		}
	}
	return n
}

func ToFloatByRemoveString(str string) (float64, error) {
	strNumber := ""
	slice := strings.Split(str, "")
	for _, s := range slice {
		if s == "." || s == "-" {
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
	fmt.Printf("%v,%v,%v,%v,%v,%v\n",
		GetSecuritiesCompanyName(stock.SecuritiesCompany),
		GetCountryName(stock.StockCountry),
		GetSecuritiesAccountName(stock.SecuritiesAccount),
		stock.SecuritiesCode,
		stock.NumberOfOwnedStock,
		stock.AveragePurchasePriceOne,
	)
}

func outputPortfolioCsvFormatOne(stock Stock) {
	// コード,市場,名称,業種,保有数
	// ,購入価格(1株あたり),購入価格(合計),時価,損益(金額),損益(割合)
	// ,EPS,1株配当,配当利回り,保有数購入価格備考
	fmt.Printf(
		"%v,,%v,%v,%v"+
			",%v,%v,%v,%v,%v"+
			",,%v,%v,\n",

		stock.SecuritiesCode, // コード
		// 市場-不要
		stock.CompanyName,        // 名称
		stock.Industry,           // 業種
		stock.NumberOfOwnedStock, // 保有数

		stock.AveragePurchasePriceOne,                     // 購入価格(1株あたり)
		stock.AveragePurchasePriceAll,                     // 購入価格(合計)
		stock.ValuationAll,                                // 時価(合計)
		stock.ProfitAndLossAll,                            // 損益(金額)
		fmt.Sprintf("%.2f", stock.ProfitAndLossRatio)+"%", // 損益(割合)

		// EPS-不要
		fmt.Sprintf("%.2f", stock.DividendOne)+"円",   // 1株配当
		fmt.Sprintf("%.2f", stock.DividendRatio)+"%", // 配当利回り
		// 保有数購入価格備考-不要
	)
}

func OutputPortfolioCsvFormatAll(stocks []Stock) {
	for _, stock := range stocks {
		outputPortfolioCsvFormatOne(stock)
	}
}

func GetStringOnlyInsideBrackets(str string) string {
	reg1 := regexp.MustCompile(`(.+)\(`)
	reg2 := regexp.MustCompile(`\)(.+)`)
	str = reg1.ReplaceAllString(str, "")
	str = reg2.ReplaceAllString(str, "")
	return str
}
