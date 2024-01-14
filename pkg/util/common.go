package util

import (
	"fmt"
	"github.com/okinari/golibs"
	"regexp"
	"strings"
)

type SecuritiesCompany int

const WaitTimeSecond = 3

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
	NisaOldAccount
	NisaNewAccount
	SpecificAccount
	GeneralAccount
)

const (
	// Yahooファイナンスのポートフォリオ画面のID(ポートフォリオID)を取得する
	// 証券の情報によって分類している
	// 分類は以下の通り
	// 楽天(日本株、外国株):3
	// SBI(特定)(1000-9000番台):8-16
	// SBI(旧NISA)(1000-9000番台):17-25
	// SBI(外国株):26
	// YahooFinance
	YahooPortfolioIDRakuten            = "3"
	YahooPortfolioIDSBISpecific1000    = "8"
	YahooPortfolioIDSBISpecific2000    = "9"
	YahooPortfolioIDSBISpecific3000    = "10"
	YahooPortfolioIDSBISpecific4000    = "11"
	YahooPortfolioIDSBISpecific5000    = "12"
	YahooPortfolioIDSBISpecific6000    = "13"
	YahooPortfolioIDSBISpecific7000    = "14"
	YahooPortfolioIDSBISpecific8000    = "15"
	YahooPortfolioIDSBISpecific9000    = "16"
	YahooPortfolioIDSBINisaOld1000     = "17"
	YahooPortfolioIDSBINisaOld2000     = "18"
	YahooPortfolioIDSBINisaOld3000     = "19"
	YahooPortfolioIDSBINisaOld4000     = "20"
	YahooPortfolioIDSBINisaOld5000     = "21"
	YahooPortfolioIDSBINisaOld6000     = "22"
	YahooPortfolioIDSBINisaOld7000     = "23"
	YahooPortfolioIDSBINisaOld8000     = "24"
	YahooPortfolioIDSBINisaOld9000     = "25"
	YahooPortfolioIDSBIForeignGeneral  = "26"
	YahooPortfolioIDSBIForeignSpecific = "27"
	YahooPortfolioIDSBIForeignOldNisa  = "28"
	YahooPortfolioIDSBIForeignNewNisa  = "29"
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
	EarningsPerShare        float64           // EPS(1株あたり純利益)
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
	default:
		panic("unhandled default case")
	}
}

func GetSecuritiesAccount(str string) SecuritiesAccount {
	if strings.Contains(str, "旧NISA") {
		return NisaOldAccount
	}
	if strings.Contains(str, "新NISA") {
		return NisaNewAccount
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
	case NisaOldAccount:
		return "旧NISA口座"
	case NisaNewAccount:
		return "新NISA口座"
	case SpecificAccount:
		return "特定口座"
	case GeneralAccount:
		return "一般口座"
	default:
		panic("unhandled default case")
	}
}

func GetCountry(str string) Country {
	if strings.Contains(str, "日本") {
		return Japan
	}
	if strings.Contains(str, "japan") {
		return Japan
	}
	if strings.Contains(str, "Japan") {
		return Japan
	}
	if strings.Contains(str, "米国") {
		return Usa
	}
	if strings.Contains(str, "usa") {
		return Usa
	}
	if strings.Contains(str, "USA") {
		return Usa
	}
	if strings.Contains(str, "アメリカ") {
		return Usa
	}
	if strings.Contains(str, "america") {
		return Usa
	}
	if strings.Contains(str, "America") {
		return Usa
	}
	return NoneCountry
}

func GetCountryName(country Country) string {
	switch country {
	case Japan:
		return "日本"
	case Usa:
		return "米国"
	default:
		panic("unhandled default case")
	}
}

func WaitTime() {
	golibs.WaitTimeSecond(WaitTimeSecond)
}

func ToIntByRemoveString(str string) int {
	return golibs.ToIntByRemoveString(str)
}

func ToFloatByRemoveString(str string) (float64, error) {
	return golibs.ToFloatByRemoveString(str)
}

func ToStringByFloat64(number float64) string {
	return golibs.ToStringByFloat64(number)
}

func DiffStocks(stocksMain, stocksSub []Stock) []Stock {
	var diffStocks []Stock
	for _, stockMain := range stocksMain {
		isNotExist := true
		for _, stockSub := range stocksSub {
			if DiffStock(stockMain, stockSub) {
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

func DiffStock(stockMain, stockSub Stock) bool {
	if stockSub.SecuritiesCompany == stockMain.SecuritiesCompany &&
		stockSub.StockCountry == stockMain.StockCountry &&
		stockSub.SecuritiesAccount == stockMain.SecuritiesAccount &&
		stockSub.SecuritiesCode == stockMain.SecuritiesCode &&
		stockSub.NumberOfOwnedStock == stockMain.NumberOfOwnedStock &&
		stockSub.AveragePurchasePriceOne == stockMain.AveragePurchasePriceOne {
		return true
	}

	return false
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
		"%v,%v,,%v,%v,%v"+
			",%v,%v,%v,%v,%v"+
			",,%v,%v,\n",
		//"%v,,%v,%v,%v"+
		//	",%v,%v,%v,%v,%v"+
		//	",,%v,%v,\n",

		stock.SecuritiesCompany, // 証券会社
		stock.SecuritiesCode,    // コード
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

// GetPortfolioID
// yahooファイナンスのポートフォリオ画面のID(ポートフォリオID)を取得する
// 証券の情報によって分類している
// 分類は以下の通り
// 楽天(日本株、外国株):3
// SBI(特定)(1000-9000番台):8-16
// SBI(旧NISA)(1000-9000番台):17-25
// SBI(外国株)(一般):26
// SBI(外国株)(特定):27
// SBI(外国株)(旧NISA):28
// SBI(外国株)(新NISA):29
func GetPortfolioID(stock Stock) string {

	if stock.SecuritiesCompany == GetSecuritiesCompany("楽天") {
		return YahooPortfolioIDRakuten
	}

	if stock.SecuritiesCompany == GetSecuritiesCompany("SBI") {

		if stock.StockCountry == GetCountry("米国") {
			if stock.SecuritiesAccount == GetSecuritiesAccount("一般") {
				return YahooPortfolioIDSBIForeignGeneral
			}
			if stock.SecuritiesAccount == GetSecuritiesAccount("特定") {
				return YahooPortfolioIDSBIForeignSpecific
			}
			if stock.SecuritiesAccount == GetSecuritiesAccount("旧NISA") {
				return YahooPortfolioIDSBIForeignOldNisa
			}
			if stock.SecuritiesAccount == GetSecuritiesAccount("新NISA") {
				return YahooPortfolioIDSBIForeignNewNisa
			}
		}

		if stock.StockCountry == GetCountry("日本") {
			securitiesCode := ToIntByRemoveString(stock.SecuritiesCode)

			if stock.SecuritiesAccount == GetSecuritiesAccount("特定") {
				if securitiesCode <= 1999 {
					return YahooPortfolioIDSBISpecific1000
				}
				if securitiesCode >= 2000 && securitiesCode <= 2999 {
					return YahooPortfolioIDSBISpecific2000
				}
				if securitiesCode >= 3000 && securitiesCode <= 3999 {
					return YahooPortfolioIDSBISpecific3000
				}
				if securitiesCode >= 4000 && securitiesCode <= 4999 {
					return YahooPortfolioIDSBISpecific4000
				}
				if securitiesCode >= 5000 && securitiesCode <= 5999 {
					return YahooPortfolioIDSBISpecific5000
				}
				if securitiesCode >= 6000 && securitiesCode <= 6999 {
					return YahooPortfolioIDSBISpecific6000
				}
				if securitiesCode >= 7000 && securitiesCode <= 7999 {
					return YahooPortfolioIDSBISpecific7000
				}
				if securitiesCode >= 8000 && securitiesCode <= 8999 {
					return YahooPortfolioIDSBISpecific8000
				}
				if securitiesCode >= 9000 {
					return YahooPortfolioIDSBISpecific9000
				}
			}

			if stock.SecuritiesAccount == GetSecuritiesAccount("旧NISA") {
				if securitiesCode <= 1999 {
					return YahooPortfolioIDSBINisaOld1000
				}
				if securitiesCode >= 2000 && securitiesCode <= 2999 {
					return YahooPortfolioIDSBINisaOld2000
				}
				if securitiesCode >= 3000 && securitiesCode <= 3999 {
					return YahooPortfolioIDSBINisaOld3000
				}
				if securitiesCode >= 4000 && securitiesCode <= 4999 {
					return YahooPortfolioIDSBINisaOld4000
				}
				if securitiesCode >= 5000 && securitiesCode <= 5999 {
					return YahooPortfolioIDSBINisaOld5000
				}
				if securitiesCode >= 6000 && securitiesCode <= 6999 {
					return YahooPortfolioIDSBINisaOld6000
				}
				if securitiesCode >= 7000 && securitiesCode <= 7999 {
					return YahooPortfolioIDSBINisaOld7000
				}
				if securitiesCode >= 8000 && securitiesCode <= 8999 {
					return YahooPortfolioIDSBINisaOld8000
				}
				if securitiesCode >= 9000 {
					return YahooPortfolioIDSBINisaOld9000
				}
			}
		}
	}

	return ""
}
