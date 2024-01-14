package yahoo_finance

import (
	"fmt"
	"github.com/okinari/golibs"
	"github.com/okinari/securities-portfolio/pkg/util"
	"github.com/sclevine/agouti"
)

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

	err := yf.ws.NavigatePageForWait("https://login.yahoo.co.jp/config/login?.src=finance&lg=jp&.intl=jp&.done=https%3A%2F%2Ffinance.yahoo.co.jp%2F", 3*1000, 100)
	if err != nil {
		return err
	}

	err = yf.ws.SetStringByID("login_handle", userName)
	if err != nil {
		return err
	}
	err = yf.ws.ClickButtonByText("次へ")
	if err != nil {
		return err
	}

	err = yf.ws.WaitUntilURL("https://finance.yahoo.co.jp/", 10*1000, 100)
	if err != nil {
		return err
	}

	return nil
}

func (yf *YahooFinance) GetSecuritiesAccountInfo(portfolioId string) ([]util.Stock, error) {

	var stocks []util.Stock

	err := yf.ws.NavigatePageForWait(getPortfolioPageURL(portfolioId, false), 10*1000, 100)
	if err != nil {
		return nil, err
	}

	// テーブル取得
	multiSelection := yf.ws.GetPage().First("table")

	// 列の割り出し
	stockColAndColNum := map[int]int{}
	theadColumns := multiSelection.All("thead tr th")
	for i := 0; true; i++ {
		colName, err := theadColumns.At(i).Text()
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

		tbodyColumns := row.All("td")
		stock := util.Stock{}
		for stockColumn, columnNumber := range stockColAndColNum {
			if stockColumn == None {
				continue
			}

			err = getStockInfo(stockColumn, tbodyColumns.At(columnNumber), &stock)
			if err != nil {
				return nil, err
			}
		}

		// stockの情報を詰め込む
		stocks = append(stocks, stock)
	}

	return stocks, nil
}

func getStockInfo(stockColumn int, column *agouti.Selection, stock *util.Stock) error {
	var err error

	// コード・市場・名称
	if stockColumn == CodeAndCompanyName {
		// 日本株の場合
		stock.SecuritiesCode, err = column.First("dt").Text()
		if err != nil {
			return err
		}
		// 米国株の場合
		if stock.SecuritiesCode == "" {
			stock.SecuritiesCode, err = column.All("dt").At(1).Text()
			if err != nil {
				return err
			}
		}

		stock.CompanyName, err = column.All("dd").At(1).Text()
		if err != nil {
			return err
		}
	}

	// 業種
	if stockColumn == Industry {
		stock.Industry, err = column.Text()
		if err != nil {
			return err
		}
	}

	// 保有数
	if stockColumn == NumberOfOwnedStock {
		numberOfOwnedStock, err := column.Text()
		if err != nil {
			return err
		}

		// 保有数が---の場合はスキップ
		if numberOfOwnedStock == "---" {
			return nil
		}

		stock.NumberOfOwnedStock, err = util.ToFloatByRemoveString(numberOfOwnedStock)
		if err != nil {
			return err
		}
	}

	// 購入価格
	if stockColumn == AveragePurchasePriceOne {
		averagePurchasePriceOne, err := column.Text()
		if err != nil {
			return err
		}

		// 購入価格が---の場合はスキップ
		if averagePurchasePriceOne == "---" {
			return nil
		}

		stock.AveragePurchasePriceOne, err = util.ToFloatByRemoveString(averagePurchasePriceOne)
		if err != nil {
			return err
		}
	}

	// 時価
	if stockColumn == ValuationAll {
		valuationAll, err := column.Text()
		if err != nil {
			return err
		}

		// 時価が---の場合はスキップ
		if valuationAll == "---" {
			return nil
		}

		stock.ValuationAll, err = util.ToFloatByRemoveString(valuationAll)
		if err != nil {
			return err
		}
	}

	// 損益 → 損益(合計)、損益(割合)
	if stockColumn == ProfitAndLossAll {
		profitAndLossAll, err := column.All("span span span").At(0).Text()
		if err != nil {
			return err
		}

		// 損益(合計)が---の場合はスキップ
		if profitAndLossAll == "---" {
			return nil
		}

		stock.ProfitAndLossAll, err = util.ToFloatByRemoveString(profitAndLossAll)
		if err != nil {
			return err
		}

		profitAndLossRatio, err := column.All("span span span").At(1).Text()
		if err != nil {
			return err
		}

		// 損益(割合)が---の場合はスキップ
		if profitAndLossRatio == "---" {
			return nil
		}

		stock.ProfitAndLossRatio, err = util.ToFloatByRemoveString(profitAndLossRatio)
		if err != nil {
			return err
		}
	}

	// EPS
	if stockColumn == EarningsPerShare {
		earningsPerShare, err := column.All("span span span").At(1).Text()
		if err != nil {
			return err
		}

		// EPSが---の場合はスキップ
		if earningsPerShare == "---" {
			return nil
		}

		stock.EarningsPerShare, err = util.ToFloatByRemoveString(earningsPerShare)
		if err != nil {
			return err
		}
	}

	//　1株配当
	if stockColumn == DividendOne {
		dividendOne, err := column.All("span span span").At(0).Text()
		if err != nil {
			return err
		}

		// 1株配当が---の場合はスキップ
		if dividendOne == "---" {
			return nil
		}

		stock.DividendOne, err = util.ToFloatByRemoveString(dividendOne)
		if err != nil {
			return err
		}
	}

	// 配当利回り
	if stockColumn == DividendRatio {
		dividendRatio, err := column.All("span span span").At(0).Text()
		if err != nil {
			return err
		}

		// 配当利回りが---の場合はスキップ
		if dividendRatio == "---" {
			return nil
		}

		stock.DividendRatio, err = util.ToFloatByRemoveString(dividendRatio)
		if err != nil {
			return err
		}
	}

	// 取得項目以外の項目のため、特に何もしない
	return nil
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

// AddSecuritiesCode is 手動作業も含めて、証券番号を追加するための処理である
func (yf *YahooFinance) AddSecuritiesCode(portfolioId string, stocks []util.Stock) error {

	// 画面上の証券コード一覧を取得
	screenStocks, err := yf.GetSecuritiesAccountInfo(portfolioId)
	if err != nil {
		return err
	}

	var securitiesCodes []string
	for _, stock := range screenStocks {
		securitiesCodes = append(securitiesCodes, stock.SecuritiesCode)
	}

	// 証券番号の存在を確認し、存在しない場合は追加する
	for _, stock := range stocks {

		isExist := false
		for _, securitiesCode := range securitiesCodes {
			if stock.SecuritiesCode == securitiesCode {
				isExist = true
				break
			}
		}

		// 存在しない場合、証券番号を追加する
		if isExist == false {
			err = yf.addSecuritiesCodeInSearchArea(stock)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// SetSecuritiesCode 引数に与えた証券コードを画面に設定する。なお、手動の処理があるので注意
func (yf *YahooFinance) SetSecuritiesCode(portfolioId string, stocks []util.Stock) error {

	url := getPortfolioPageURL(portfolioId, true)
	err := yf.ws.NavigatePageForWait(url, 10*1000, 100)
	if err != nil {
		return err
	}

	// テーブル取得
	multiSelection := yf.ws.GetPage().First("table")

	rows := multiSelection.All("tbody tr")

	// 保有数、購入単価、の全ての項目が空もしくは0になったことを確認するまで、5分程度待機する
	for count := 100; count > 0; count-- {

		isAllEmpty := true
		for i := 0; true; i++ {

			// 最後の行まで辿りついたとき(=行の最初の列がエラーになったとき)、処理をやめる
			row := rows.At(i)
			_, err = row.All("td").At(0).Text()
			if err != nil {
				break
			}

			// 保有数を確認
			value, err := row.All("td").At(3).First("input").Attribute("value")
			if err != nil {
				return err
			}
			if value != "" && value != "0" {
				isAllEmpty = false
				break
			}

			// 購入単価を確認
			value, err = row.All("td").At(4).First("input").Attribute("value")
			if err != nil {
				return err
			}
			if value != "" && value != "0" {
				isAllEmpty = false
				break
			}
		}

		if isAllEmpty {
			break
		}

		if count%10 == 0 {
			fmt.Printf("保有数、購入単価が空になるのを待機中...。あと%v秒\n", count*3)
		}
		util.WaitTime()
	}

	// 画面から1行ずつ確認する
	for i := 0; true; i++ {

		// 行の最初の列がエラーになったら、処理をやめる
		row := rows.At(i)
		_, err = row.All("td").At(0).Text()
		if err != nil {
			break
		}

		// 画面の証券番号を取得
		// 日本株の場合
		securityCode, err := row.All("td").At(2).First("dt").Text()
		if err != nil {
			return err
		}
		// 米国株の場合
		if securityCode == "" {
			securityCode, err = row.All("td").At(2).All("dt").At(1).Text()
			if err != nil {
				return err
			}
		}

		// 引数に渡された証券情報と合致する行を探し、該当の株情報を追加、もしくは修正する
		for _, stock := range stocks {

			// 証券番号が合致しない場合、次の行へ
			if stock.SecuritiesCode != securityCode {
				continue
			}

			// 保有数を設定
			element := row.All("td").At(3).First("input")
			err := element.Clear()
			if err != nil {
				return err
			}
			err = element.Fill(util.ToStringByFloat64(stock.NumberOfOwnedStock))
			if err != nil {
				return err
			}

			// 購入単価を設定
			element = row.All("td").At(4).First("input")
			err = element.Clear()
			if err != nil {
				return err
			}
			err = element.Fill(util.ToStringByFloat64(stock.AveragePurchasePriceOne))
			if err != nil {
				return err
			}
		}
	}

	// 手動で変更するため、完了してポートフォリオ画面に遷移するまで、待機する
	err = yf.ws.WaitForURLChange(url, 600*1000, 5*1000)
	if err != nil {
		return err
	}

	return nil
}

// TODO: 未完成
// 値を入力させても、更新することができない、、、
func (yf *YahooFinance) SetSecurities(portfolioId string, stocks []util.Stock) error {

	url := getPortfolioPageURL(portfolioId, true)
	err := yf.ws.NavigatePageForWait(url, 10*1000, 100)
	if err != nil {
		return err
	}

	// テーブル取得
	multiSelection := yf.ws.GetPage().First("table")

	rows := multiSelection.All("tbody tr")

	// TODO: 証券番号の自動追加を試みたが、うまくいかなかったので、一部手動処理を入れている
	// 証券番号の存在を確認し、存在しない場合は追加する
	for _, stock := range stocks {

		// 既に証券番号が存在する場合は、次の証券番号へ
		result, err := yf.isExistStock(rows, stock)
		if err != nil {
			return err
		}
		if result {
			continue
		}

		// 存在しない場合、証券番号を追加する
		err = yf.addSecuritiesCodeInSearchArea(stock)
		if err != nil {
			return err
		}
	}

	// TODO: 最後に更新ボタンを押しても、値が更新されない
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
			element := row.All("td").At(3).First("input")
			err := element.Clear()
			if err != nil {
				return err
			}
			err = element.Fill(util.ToStringByFloat64(stock.NumberOfOwnedStock))
			if err != nil {
				return err
			}

			// 購入単価を設定
			element = row.All("td").At(4).First("input")
			err = element.Clear()
			if err != nil {
				return err
			}
			err = element.Fill(util.ToStringByFloat64(stock.AveragePurchasePriceOne))
			if err != nil {
				return err
			}
		}

		err = yf.ws.WaitUntilURL(getPortfolioPageURL(portfolioId, false), 600*1000, 5*1000)
		if err != nil {
			return err
		}
	}

	return nil
}

func (yf *YahooFinance) isExistStock(rows *agouti.MultiSelection, stock util.Stock) (bool, error) {
	// 画面から1行ずつ確認する
	for i := 0; true; i++ {

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

		// 証券番号が存在する場合、trueを返却
		if securityCode == stock.SecuritiesCode {
			return true, nil
		}
	}

	// 証券番号が存在しない場合、falseを返却
	return false, nil
}

// Yahooファイナンスの編集画面にて、検索ボックスに証券番号を追加する
func (yf *YahooFinance) addSecuritiesCodeInSearchArea(stock util.Stock) error {
	page := yf.ws.GetPage()
	input := page.All("input[type='search'][aria-label='銘柄検索']").At(0)
	if stock.StockCountry == util.GetCountry("日本") {
		stock.SecuritiesCode = stock.SecuritiesCode + ".T"
	}
	err := input.Fill(stock.SecuritiesCode)
	if err != nil {
		return err
	}

	golibs.WaitTimeMillSecond(1000)

	// 銘柄追加のボタンを押下
	err = page.All("button span").At(23).Click()
	if err != nil {
		return err
	}

	// 自動でボタンを押しても、銘柄追加が不可能なケースがあるので、その場合、手動で実行することにした
	// 銘柄コードの入力欄が空になったことを確認するまで、最大30秒程度待機する
	conditionFunction := func(limitTimeMillSecond int) (bool, error) {
		// 銘柄コード入力欄を確認
		value, err := input.Attribute("value")
		if err != nil {
			return false, err
		}
		if value == "" {
			return true, nil
		}

		if limitTimeMillSecond < 30*1000 && limitTimeMillSecond%(5*1000) == 0 {
			fmt.Printf("銘柄コードが入力完了になるのを待機中...。あと%v秒\n", limitTimeMillSecond/1000)
		}
		return false, nil
	}
	err = golibs.WaitForCondition(30*1000, 500, conditionFunction)
	if err != nil {
		return err
	}

	return nil
}

// SettingPortfolio 特定のポートフォリオの設定を変更する
func (yf *YahooFinance) SettingPortfolio(portfolioId string) error {

	err := yf.ws.NavigatePageForWait(getPortfolioPageURL(portfolioId, false), 10*1000, 100)
	if err != nil {
		return err
	}

	page := yf.ws.GetPage()
	err = page.All("button").At(4).Click()
	if err != nil {
		return err
	}

	checkedList := page.All("input[type=checkbox]")
	list := []int{0, 24, 25, 26, 28, 32, 33, 36, 46}
	for _, v := range list {
		err := checkedList.At(v).Check()
		if err != nil {
			return err
		}
	}

	err = page.All("button span span").At(1).Click()
	if err != nil {
		return err
	}

	golibs.WaitTimeSecond(3)

	return nil
}

// DeletePortfolio 特定のポートフォリオの内容を全て削除する
func (yf *YahooFinance) DeletePortfolio(portfolioId string) error {

	url := getPortfolioPageURL(portfolioId, false)
	err := yf.ws.NavigatePageForWait(url, 10*1000, 100)
	if err != nil {
		return err
	}

	page := yf.ws.GetPage()
	count, err := page.All("button[aria-label='削除']").Count()
	if err != nil {
		return err
	}

	for i := 10; i < count; i++ {
		err = page.All("button[aria-label='削除']").At(10).Click()
		if err != nil {
			return err
		}

		// 連続でボタンを押すと落ちるので、500ミリ秒待機する
		golibs.WaitTimeMillSecond(500)

		err = page.All("button span span").At(4).Click()
		if err != nil {
			return err
		}

		// 連続でボタンを押すと落ちるので、500ミリ秒待機する
		golibs.WaitTimeMillSecond(500)
	}

	return nil
}

func getPortfolioPageURL(portfolioID string, isEdit bool) string {
	if isEdit {
		return "https://finance.yahoo.co.jp/portfolio/detail/edit?portfolioId=" + portfolioID
	} else {
		return "https://finance.yahoo.co.jp/portfolio/detail?portfolioId=" + portfolioID
	}
}
