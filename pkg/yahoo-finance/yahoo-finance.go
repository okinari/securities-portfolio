package yahoo_finance

import (
	"fmt"
	"github.com/okinari/golibs"
	"github.com/okinari/securities-portfolio/pkg/util"
	"github.com/sclevine/agouti"
)

const LoginUrl = "https://login.yahoo.co.jp/config/login?.src=finance&lg=jp&.intl=jp&.done=https%3A%2F%2Ffinance.yahoo.co.jp%2F"

type YahooFinance struct {
	ws *golibs.WebScraping
}

func NewYahooFinance() (*YahooFinance, error) {
	ws, err := golibs.NewWebScraping(golibs.IsHeadless(false))
	if err != nil {
		return nil, err
	}

	return &YahooFinance{
		ws: ws,
	}, nil
}

func (yf *YahooFinance) Close() error {
	err := yf.ws.Close()
	if err != nil {
		return err
	}
	return nil
}

func (yf *YahooFinance) Login(userName string) error {

	err := yf.ws.NavigatePage(LoginUrl)
	if err != nil {
		return err
	}
	util.WaitTime()

	err = yf.ws.SetStringByID("login_handle", userName)
	if err != nil {
		return err
	}
	err = yf.ws.ClickButtonByText("次へ")
	if err != nil {
		return err
	}

	page := yf.ws.GetPage()
	count := 20
	for count > 0 {
		util.WaitTime()
		url, err := page.URL()
		if err != nil {
			return err
		}
		if url == "https://finance.yahoo.co.jp/" {
			break
		}
		count--
	}

	return nil
}

func (yf *YahooFinance) GetSecuritiesAccountInfo(portfolioId string) ([]util.Stock, error) {

	var stocks []util.Stock

	err := yf.ws.NavigatePage("https://finance.yahoo.co.jp/portfolio/detail?portfolioId=" + portfolioId)
	if err != nil {
		return nil, err
	}
	util.WaitTime()

	// テーブル取得
	multiSelection := yf.ws.GetPage().First("table")

	// 列の割り出し
	stockColAndColNum := map[int]int{}
	columns := multiSelection.All("thead tr th")
	for i := 0; true; i++ {
		colName, err := columns.At(i).Text()
		if err != nil {
			break
		}
		stockColAndColNum[getColumnNameMapping(colName)] = i
	}

	stocks = []util.Stock{}

	// データの抽出
	rows := multiSelection.All("tbody tr")
	for i := 0; true; i++ {

		// 最初の列がエラーになったら、処理をやめる
		row := rows.At(i)
		_, err = row.All("td").At(0).Text()
		if err != nil {
			break
		}

		stock, err := setStock(stockColAndColNum, row)
		if err != nil {
			return nil, err
		}

		// stockの情報を詰め込む
		stocks = append(stocks, *stock)
	}

	return stocks, nil
}

func setStock(stockColAndColNum map[int]int, row *agouti.Selection) (*util.Stock, error) {
	var err error
	stock := util.Stock{}

	for stockColumn, columnNumber := range stockColAndColNum {
		if stockColumn == None {
			continue
		}

		// コード・市場・名称
		if stockColumn == CodeAndCompanyName {
			stock.SecuritiesCode, err = row.All("td").At(columnNumber).First("dt").Text()
			if err != nil {
				return nil, err
			}

			stock.CompanyName, err = row.All("td").At(columnNumber).All("dd").At(1).Text()
			if err != nil {
				return nil, err
			}
		}

		// 業種
		if stockColumn == Industry {
			stock.Industry, err = row.All("td").At(columnNumber).Text()
			if err != nil {
				return nil, err
			}
		}

		// 保有数
		if stockColumn == NumberOfOwnedStock {
			numberOfOwnedStock, err := row.All("td").At(columnNumber).Text()
			if err != nil {
				return nil, err
			}

			// 保有数が---の場合はスキップ
			if numberOfOwnedStock == "---" {
				continue
			}

			stock.NumberOfOwnedStock, err = util.ToFloatByRemoveString(numberOfOwnedStock)
			if err != nil {
				return nil, err
			}
		}

		// 購入価格
		if stockColumn == AveragePurchasePriceOne {
			averagePurchasePriceOne, err := row.All("td").At(columnNumber).Text()
			if err != nil {
				return nil, err
			}

			// 購入価格が---の場合はスキップ
			if averagePurchasePriceOne == "---" {
				continue
			}

			stock.AveragePurchasePriceOne, err = util.ToFloatByRemoveString(averagePurchasePriceOne)
			if err != nil {
				return nil, err
			}
		}

		// 時価
		if stockColumn == ValuationAll {
			valuationAll, err := row.All("td").At(columnNumber).Text()
			if err != nil {
				return nil, err
			}

			// 時価が---の場合はスキップ
			if valuationAll == "---" {
				continue
			}

			stock.ValuationAll, err = util.ToFloatByRemoveString(valuationAll)
			if err != nil {
				return nil, err
			}
		}

		// 損益 → 損益(合計)、損益(割合)
		if stockColumn == ProfitAndLossAll {
			profitAndLossAll, err := row.All("td").At(columnNumber).All("span span span").At(0).Text()
			if err != nil {
				return nil, err
			}

			// 損益が---の場合はスキップ
			if profitAndLossAll == "---" {
				continue
			}

			stock.ProfitAndLossAll, err = util.ToFloatByRemoveString(profitAndLossAll)
			if err != nil {
				return nil, err
			}

			profitAndLossRatio, err := row.All("td").At(columnNumber).All("span span span").At(1).Text()
			if err != nil {
				return nil, err
			}
			stock.ProfitAndLossRatio, err = util.ToFloatByRemoveString(profitAndLossRatio)
			if err != nil {
				return nil, err
			}
		}

		// EPS
		if stockColumn == EarningsPerShare {
			earningsPerShare, err := row.All("td").At(columnNumber).All("span span span").At(1).Text()
			if err != nil {
				return nil, err
			}

			// EPSが---の場合はスキップ
			if earningsPerShare == "---" {
				continue
			}

			stock.EarningsPerShare, err = util.ToFloatByRemoveString(earningsPerShare)
			if err != nil {
				return nil, err
			}
		}

		//　1株配当
		if stockColumn == DividendOne {
			dividendOne, err := row.All("td").At(columnNumber).All("span span span").At(0).Text()
			if err != nil {
				return nil, err
			}

			// 1株配当が---の場合はスキップ
			if dividendOne == "---" {
				continue
			}

			stock.DividendOne, err = util.ToFloatByRemoveString(dividendOne)
			if err != nil {
				return nil, err
			}
		}

		// 配当利回り
		if stockColumn == DividendRatio {
			dividendRatio, err := row.All("td").At(columnNumber).All("span span span").At(0).Text()
			if err != nil {
				return nil, err
			}

			// 配当利回りが---の場合はスキップ
			if dividendRatio == "---" {
				continue
			}

			stock.DividendRatio, err = util.ToFloatByRemoveString(dividendRatio)
			if err != nil {
				return nil, err
			}
		}
	}

	return &stock, nil
}

const (
	None = iota
	CodeAndCompanyName
	Industry
	NumberOfOwnedStock
	AveragePurchasePriceOne
	ValuationAll
	ProfitAndLossAll
	EarningsPerShare
	DividendOne
	DividendRatio
)

func getColumnNameMapping(columnName string) int {
	switch columnName {
	case "コード・市場・名称":
		return CodeAndCompanyName
	case "業種":
		return Industry
	case "保有数":
		return NumberOfOwnedStock
	case "購入価格":
		return AveragePurchasePriceOne
	case "時価":
		return ValuationAll
	case "損益":
		return ProfitAndLossAll
	case "EPS":
		return EarningsPerShare
	case "1株配当":
		return DividendOne
	case "配当利回り":
		return DividendRatio
	default:
		return None
	}
}

// TODO: 未完成
// 値を入力させても、更新することができない、、、
func (yf *YahooFinance) SetSecurities(portfolioId string, stocks []util.Stock) error {

	fmt.Println("start")

	err := yf.ws.NavigatePage("https://finance.yahoo.co.jp/portfolio/detail/edit?portfolioId=" + portfolioId)
	if err != nil {
		return err
	}
	util.WaitTime()

	fmt.Println("navigate ok")

	// テーブル取得
	multiSelection := yf.ws.GetPage().First("table")

	fmt.Println("code add start")

	rows := multiSelection.All("tbody tr")

	// TODO: 証券番号の自動追加を試みたが、うまくいかなかったので、一旦コメントアウト
	//// 証券番号の存在を確認し、存在しない場合は追加する
	//for _, stock := range stocks {
	//
	//	fmt.Printf("code: %v\n", stock.SecuritiesCode)
	//
	//	result, err := yf.isExistStock(rows, stock)
	//	if err != nil {
	//		return err
	//	}
	//
	//	if result == false {
	//		page := yf.ws.GetPage()
	//		input := page.All("input[type='search'][aria-label='銘柄検索']").At(0)
	//		err = input.Fill(stock.SecuritiesCode)
	//		if err != nil {
	//			fmt.Printf("no code error: %v\n", stock.SecuritiesCode)
	//			return err
	//		}
	//		err = input.Fill(".")
	//		if err != nil {
	//			fmt.Printf("no code error: %v\n", stock.SecuritiesCode)
	//			return err
	//		}
	//		err = input.Fill("T")
	//		if err != nil {
	//			fmt.Printf("no code error: %v\n", stock.SecuritiesCode)
	//			return err
	//		}
	//		err = page.FindByButton("銘柄追加").Click()
	//		if err != nil {
	//			fmt.Printf("no 銘柄追加 error: %v\n", stock.SecuritiesCode)
	//			return err
	//		}
	//	}
	//
	//	util.WaitTime()
	//}

	// 画面から1行ずつ確認する
	for i := 0; true; i++ {

		// 行の最初の列がエラーになったら、処理をやめる
		row := rows.At(i)
		_, err = row.All("td").At(0).Text()
		if err != nil {
			break
		}

		// 画面の証券番号を取得
		securityCode, err := row.All("td").At(2).First("dt").Text()
		if err != nil {
			return err
		}

		// 引数に渡された証券情報と合致する行を探し、該当の株情報を追加、もしくは修正する
		for _, stock := range stocks {

			// 証券番号が合致しない場合、次の行へ
			if stock.SecuritiesCode != securityCode {
				continue
			}

			// 保有数を設定
			//element := row.All("td").At(3).First("input")
			//err := element.Clear()
			//if err != nil {
			//	return err
			//}
			//err = element.Fill(util.ToStringByFloat64(stock.NumberOfOwnedStock))
			//if err != nil {
			//	return err
			//}
			err = yf.ws.ExecJavaScript("document.querySelector('table').querySelectorAll('tbody tr')["+util.ToStringByInt(i)+"].querySelectorAll('td')[3].querySelector('input').value = '"+util.ToStringByFloat64(stock.NumberOfOwnedStock)+"'", nil)
			if err != nil {
				return err
			}

			// 購入単価を設定
			//element = row.All("td").At(4).First("input")
			//err = element.Clear()
			//if err != nil {
			//	return err
			//}
			//err = element.Fill(util.ToStringByFloat64(stock.AveragePurchasePriceOne))
			//if err != nil {
			//	return err
			//}
			err = yf.ws.ExecJavaScript("document.querySelector('table').querySelectorAll('tbody tr')["+util.ToStringByInt(i)+"].querySelectorAll('td')[4].querySelector('input').value = '"+util.ToStringByFloat64(stock.AveragePurchasePriceOne)+"'", nil)
			if err != nil {
				return err
			}
		}
	}

	for i := 20; i > 0; i-- {
		util.WaitTime()
	}

	return nil
}

func isNotExistStock(securitiesCode string, stocks []util.Stock) bool {
	for _, stock := range stocks {
		if stock.SecuritiesCode == securitiesCode {
			return false
		}
	}
	return true
}

func (yf *YahooFinance) isExistStock(rows *agouti.MultiSelection, stock util.Stock) (bool, error) {
	// 画面から1行ずつ確認する
	for i := 0; true; i++ {

		fmt.Printf("loop: %v\n", i)

		// 行の最初の列がエラーになったら、処理をやめる
		row := rows.At(i)
		_, err := row.All("td").At(0).Text()
		if err != nil {
			break
		}

		// 画面の証券番号を取得
		securityCode, err := row.All("td").At(2).First("dt").Text()
		if err != nil {
			return false, err
		}

		// 証券番号が存在しない場合、追加する
		if securityCode == stock.SecuritiesCode {
			return true, nil
		}
	}

	return false, nil
}
