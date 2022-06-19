package rakuten_securities

import (
	"github.com/okinari/golibs"
	"github.com/okinari/securities-portfolio/pkg/util"
	"strconv"
)

const LoginUrl = "https://www.rakuten-sec.co.jp/ITS/V_ACT_Login.html"

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

	err := rs.ws.NavigatePage(LoginUrl)
	if err != nil {
		return err
	}

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

func (rs *RakutenSecurities) GetSecuritiesAccountInfo() ([]util.StockInfo, error) {

	var stockInfoList []util.StockInfo

	// 国内株式（現物/特定預り）
	siList, err := rs.getStockListJapanForJapanSpecificAccount()
	if err != nil {
		return nil, err
	}
	stockInfoList = append(stockInfoList, siList...)

	// 国内株式（現物/NISA預り）
	siList, err = rs.getStockListForJapanNisaAccount()
	if err != nil {
		return nil, err
	}
	stockInfoList = append(stockInfoList, siList...)

	// 米国株式（現物/特定預り）
	siList, err = rs.getStockListForUsSpecificAccount()
	if err != nil {
		return nil, err
	}
	stockInfoList = append(stockInfoList, siList...)

	// 米国株式（現物/NISA預り）
	siList, err = rs.getStockListForUsNisaAccount()
	if err != nil {
		return nil, err
	}
	stockInfoList = append(stockInfoList, siList...)

	return stockInfoList, nil
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

func (rs *RakutenSecurities) getStockListJapanForJapanSpecificAccount() ([]util.StockInfo, error) {

	err := rs.openJapanStockScreen()
	if err != nil {
		return nil, err
	}

	multiSelection := rs.ws.GetPage().All("table#poss-tbl-sp > tbody > tr[align=right]")

	count, _ := multiSelection.Count()
	stockInfoList := make([]util.StockInfo, count)

	arrayCount := 0
	for i := 0; ; i++ {

		stockInfo := &util.StockInfo{
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
		stockInfo.SecuritiesCode = strconv.Itoa(util.ToIntByRemoveString(secCode))

		// 保有件数
		numOfStock, err := ms.At(2).Find("a").Text()
		if err != nil {
			break
		}
		stockInfo.NumberOfOwnedStock = util.ToIntByRemoveString(numOfStock)

		// 取得単価
		priceOfAvg, err := ms.At(5).Text()
		if err != nil {
			break
		}
		stockInfo.AveragePurchasePrice, err = util.ToFloatByRemoveString(priceOfAvg)
		if err != nil {
			break
		}

		stockInfoList[arrayCount] = *stockInfo
		arrayCount++
	}

	return stockInfoList, nil
}

func (rs *RakutenSecurities) getStockListForJapanNisaAccount() ([]util.StockInfo, error) {

	err := rs.openJapanStockScreen()
	if err != nil {
		return nil, err
	}

	multiSelection := rs.ws.GetPage().All("table#poss-tbl-nisa > tbody > tr[align=right]")

	count, _ := multiSelection.Count()
	stockInfoList := make([]util.StockInfo, count)

	arrayCount := 0
	for i := 0; ; i++ {

		stockInfo := &util.StockInfo{
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
		stockInfo.SecuritiesCode = strconv.Itoa(util.ToIntByRemoveString(secCode))

		// 保有件数
		numOfStock, err := ms.At(1).Text()
		if err != nil {
			break
		}
		stockInfo.NumberOfOwnedStock = util.ToIntByRemoveString(numOfStock)

		// 取得単価
		priceOfAvg, err := ms.At(2).Text()
		if err != nil {
			break
		}
		stockInfo.AveragePurchasePrice, err = util.ToFloatByRemoveString(priceOfAvg)
		if err != nil {
			break
		}

		stockInfoList[arrayCount] = *stockInfo
		arrayCount++
	}

	return stockInfoList, nil
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

func (rs *RakutenSecurities) getStockListForUsSpecificAccount() ([]util.StockInfo, error) {

	err := rs.openUsStockScreen()
	if err != nil {
		return nil, err
	}

	multiSelection := rs.ws.GetPage().All("table.ex-jpx-change-tbl > tbody").At(0).All("tr")

	count, _ := multiSelection.Count()
	count = (count - 1) / 4
	stockInfoList := make([]util.StockInfo, count)

	arrayCount := 0
	for i := 1; ; i = i + 4 {

		stockInfo := &util.StockInfo{
			SecuritiesCompany: util.RakutenSecurities,
			StockCountry:      util.America,
			SecuritiesAccount: util.SpecificAccount,
		}

		ms := multiSelection.At(i).All("td > div > nobr")

		// 証券コード
		secCode, err := ms.At(0).Text()
		if err != nil {
			break
		}
		stockInfo.SecuritiesCode = secCode

		// 保有件数
		numOfStock, err := ms.At(1).Text()
		if err != nil {
			break
		}
		stockInfo.NumberOfOwnedStock = util.ToIntByRemoveString(numOfStock)

		// 取得単価
		priceOfAvg, err := ms.At(2).Text()
		if err != nil {
			break
		}
		stockInfo.AveragePurchasePrice, err = util.ToFloatByRemoveString(priceOfAvg)
		if err != nil {
			break
		}

		stockInfoList[arrayCount] = *stockInfo
		arrayCount++
	}

	return stockInfoList, nil
}

func (rs *RakutenSecurities) getStockListForUsNisaAccount() ([]util.StockInfo, error) {

	err := rs.openUsStockScreen()
	if err != nil {
		return nil, err
	}

	multiSelection := rs.ws.GetPage().All("table.ex-jpx-change-tbl > tbody").At(1).All("tr")

	count, _ := multiSelection.Count()
	count = (count - 1) / 4
	stockInfoList := make([]util.StockInfo, count)

	arrayCount := 0
	for i := 1; ; i = i + 4 {

		stockInfo := &util.StockInfo{
			SecuritiesCompany: util.RakutenSecurities,
			StockCountry:      util.America,
			SecuritiesAccount: util.NisaAccount,
		}

		sel := multiSelection.At(i)

		// 証券コード
		secCode, err := sel.Find("td > table > tbody >tr > td > div").Text()
		if err != nil {
			break
		}
		stockInfo.SecuritiesCode = secCode

		ms := sel.All("td > div > nobr")

		// 保有件数
		numOfStock, err := ms.At(0).Text()
		if err != nil {
			break
		}
		stockInfo.NumberOfOwnedStock = util.ToIntByRemoveString(numOfStock)

		// 取得単価
		priceOfAvg, err := ms.At(1).Text()
		if err != nil {
			break
		}
		stockInfo.AveragePurchasePrice, err = util.ToFloatByRemoveString(priceOfAvg)
		if err != nil {
			break
		}

		stockInfoList[arrayCount] = *stockInfo
		arrayCount++
	}

	return stockInfoList, nil
}
