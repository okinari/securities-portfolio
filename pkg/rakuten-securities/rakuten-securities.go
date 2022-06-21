package rakuten_securities

import (
	"github.com/okinari/golibs"
	"github.com/okinari/securities-portfolio/pkg/util"
	"strconv"
)

type RakutenSecurities struct {
	ws *golibs.WebScraping
}

func NewRakutenSecurities() (*RakutenSecurities, error) {
	ws, err := golibs.NewWebScraping(golibs.IsHeadless(false))
	if err != nil {
		return nil, err
	}

	return &RakutenSecurities{
		ws: ws,
	}, nil
}

func (rs *RakutenSecurities) Close() error {
	err := rs.ws.Close()
	if err != nil {
		return err
	}
	return nil
}

func (rs *RakutenSecurities) Login(userName string, password string) error {

	err := rs.ws.NavigatePage("https://www.rakuten-sec.co.jp/ITS/V_ACT_Login.html")
	if err != nil {
		return err
	}
	util.WaitTime()

	err = rs.ws.SetStringByID("form-login-id", userName)
	if err != nil {
		return err
	}
	err = rs.ws.SetStringByID("form-login-pass", password)
	if err != nil {
		return err
	}
	err = rs.ws.ClickButtonByID("login-btn")
	if err != nil {
		return err
	}

	util.WaitTime()

	return nil
}

func (rs *RakutenSecurities) GetSecuritiesAccountInfo() ([]util.Stock, error) {

	var stocks []util.Stock

	// 国内株式（現物/特定預り）
	tmpStocks, err := rs.GetStocksJapanForJapanSpecificAccount()
	if err != nil {
		return nil, err
	}
	stocks = append(stocks, tmpStocks...)

	// 国内株式（現物/NISA預り）
	tmpStocks, err = rs.GetStocksForJapanNisaAccount()
	if err != nil {
		return nil, err
	}
	stocks = append(stocks, tmpStocks...)

	// 米国株式（現物/特定預り）
	tmpStocks, err = rs.GetStocksForUsSpecificAccount()
	if err != nil {
		return nil, err
	}
	stocks = append(stocks, tmpStocks...)

	// 米国株式（現物/NISA預り）
	tmpStocks, err = rs.GetStocksForUsNisaAccount()
	if err != nil {
		return nil, err
	}
	stocks = append(stocks, tmpStocks...)

	return stocks, nil
}

func (rs *RakutenSecurities) openJapanStockScreen() error {

	// 国内株式の画面を開く
	err := rs.ws.ExecJavaScript("document.querySelector('a[data-ratid=mem_pc_mymenu_jp-possess-lst]').click()", nil)
	if err != nil {
		return err
	}

	util.WaitTime()

	return nil
}

func (rs *RakutenSecurities) GetStocksJapanForJapanSpecificAccount() ([]util.Stock, error) {

	err := rs.openJapanStockScreen()
	if err != nil {
		return nil, err
	}

	multiSelection := rs.ws.GetPage().All("table#poss-tbl-sp > tbody > tr[align=right]")

	count, _ := multiSelection.Count()
	stocks := make([]util.Stock, count)

	arrayCount := 0
	for i := 0; ; i++ {

		stock := &util.Stock{
			SecuritiesCompany: util.RakutenSecurities,
			StockCountry:      util.Japan,
			SecuritiesAccount: util.SpecificAccount,
		}

		ms := multiSelection.At(i).All("td td nobr")

		// 証券コード
		secCode, err := ms.At(0).Text()
		if err != nil {
			break
		}
		stock.SecuritiesCode = strconv.Itoa(util.ToIntByRemoveString(secCode))

		// 保有件数
		numOfStock, err := ms.At(2).Find("a").Text()
		if err != nil {
			break
		}
		stock.NumberOfOwnedStock = util.ToIntByRemoveString(numOfStock)

		// 取得単価
		priceOfAvg, err := ms.At(5).Text()
		if err != nil {
			break
		}
		stock.AveragePurchasePrice, err = util.ToFloatByRemoveString(priceOfAvg)
		if err != nil {
			break
		}

		stocks[arrayCount] = *stock
		arrayCount++
	}

	return stocks, nil
}

func (rs *RakutenSecurities) GetStocksForJapanNisaAccount() ([]util.Stock, error) {

	err := rs.openJapanStockScreen()
	if err != nil {
		return nil, err
	}

	multiSelection := rs.ws.GetPage().All("table#poss-tbl-nisa > tbody > tr[align=right]")

	count, _ := multiSelection.Count()
	stocks := make([]util.Stock, count)

	arrayCount := 0
	for i := 0; ; i++ {

		stock := &util.Stock{
			SecuritiesCompany: util.RakutenSecurities,
			StockCountry:      util.Japan,
			SecuritiesAccount: util.NisaAccount,
		}

		ms := multiSelection.At(i).All("td td nobr")

		// 証券コード
		secCode, err := ms.At(0).Text()
		if err != nil {
			break
		}
		stock.SecuritiesCode = strconv.Itoa(util.ToIntByRemoveString(secCode))

		// 保有件数
		numOfStock, err := ms.At(1).Text()
		if err != nil {
			break
		}
		stock.NumberOfOwnedStock = util.ToIntByRemoveString(numOfStock)

		// 取得単価
		priceOfAvg, err := ms.At(2).Text()
		if err != nil {
			break
		}
		stock.AveragePurchasePrice, err = util.ToFloatByRemoveString(priceOfAvg)
		if err != nil {
			break
		}

		stocks[arrayCount] = *stock
		arrayCount++
	}

	return stocks, nil
}

func (rs *RakutenSecurities) openUsStockScreen() error {

	// 米国株式の画面を開く
	err := rs.ws.ExecJavaScript("document.querySelector('a[data-ratid=mem_pc_mymenu_us-possess-lst]').click()", nil)
	if err != nil {
		return err
	}

	util.WaitTime()

	return nil
}

func (rs *RakutenSecurities) GetStocksForUsSpecificAccount() ([]util.Stock, error) {

	err := rs.openUsStockScreen()
	if err != nil {
		return nil, err
	}

	multiSelection := rs.ws.GetPage().All("table.ex-jpx-change-tbl > tbody").At(0).All("tr")

	count, _ := multiSelection.Count()
	count = (count - 1) / 4
	stocks := make([]util.Stock, count)

	arrayCount := 0
	for i := 1; ; i = i + 4 {

		stock := &util.Stock{
			SecuritiesCompany: util.RakutenSecurities,
			StockCountry:      util.Usa,
			SecuritiesAccount: util.SpecificAccount,
		}

		ms := multiSelection.At(i).All("td > div > nobr")

		// 証券コード
		secCode, err := ms.At(0).Text()
		if err != nil {
			break
		}
		stock.SecuritiesCode = secCode

		// 保有件数
		numOfStock, err := ms.At(1).Text()
		if err != nil {
			break
		}
		stock.NumberOfOwnedStock = util.ToIntByRemoveString(numOfStock)

		// 取得単価
		priceOfAvg, err := ms.At(2).Text()
		if err != nil {
			break
		}
		stock.AveragePurchasePrice, err = util.ToFloatByRemoveString(priceOfAvg)
		if err != nil {
			break
		}

		stocks[arrayCount] = *stock
		arrayCount++
	}

	return stocks, nil
}

func (rs *RakutenSecurities) GetStocksForUsNisaAccount() ([]util.Stock, error) {

	err := rs.openUsStockScreen()
	if err != nil {
		return nil, err
	}

	multiSelection := rs.ws.GetPage().All("table.ex-jpx-change-tbl > tbody").At(1).All("tr")

	count, _ := multiSelection.Count()
	count = (count - 1) / 4
	stocks := make([]util.Stock, count)

	arrayCount := 0
	for i := 1; ; i = i + 4 {

		stock := &util.Stock{
			SecuritiesCompany: util.RakutenSecurities,
			StockCountry:      util.Usa,
			SecuritiesAccount: util.NisaAccount,
		}

		sel := multiSelection.At(i)

		// 証券コード
		secCode, err := sel.Find("td > table > tbody >tr > td > div").Text()
		if err != nil {
			break
		}
		stock.SecuritiesCode = secCode

		ms := sel.All("td > div > nobr")

		// 保有件数
		numOfStock, err := ms.At(0).Text()
		if err != nil {
			break
		}
		stock.NumberOfOwnedStock = util.ToIntByRemoveString(numOfStock)

		// 取得単価
		priceOfAvg, err := ms.At(1).Text()
		if err != nil {
			break
		}
		stock.AveragePurchasePrice, err = util.ToFloatByRemoveString(priceOfAvg)
		if err != nil {
			break
		}

		stocks[arrayCount] = *stock
		arrayCount++
	}

	return stocks, nil
}
