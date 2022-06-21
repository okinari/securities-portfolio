package sbi_securities

import (
	"fmt"
	"github.com/okinari/golibs"
	"github.com/okinari/securities-portfolio/pkg/util"
	"github.com/sclevine/agouti"
	"strconv"
)

type SbiSecurities struct {
	ws *golibs.WebScraping
}

func NewSbiSecurities() (*SbiSecurities, error) {
	ws, err := golibs.NewWebScraping(golibs.IsHeadless(false))
	if err != nil {
		return nil, err
	}

	return &SbiSecurities{
		ws: ws,
	}, nil
}

func (ss *SbiSecurities) Close() error {
	err := ss.ws.Close()
	if err != nil {
		return err
	}
	return nil
}

func (ss *SbiSecurities) Login(userName string, password string) error {

	err := ss.ws.NavigatePage("https://www.sbisec.co.jp/ETGate")
	if err != nil {
		return err
	}

	err = ss.ws.SetStringByName("user_id", userName)
	if err != nil {
		return err
	}
	err = ss.ws.SetStringByName("user_password", password)
	if err != nil {
		return err
	}
	err = ss.ws.ClickButtonByName("ACT_login")
	if err != nil {
		return err
	}

	util.WaitTime()

	return nil
}

func (ss *SbiSecurities) GetSecuritiesAccountInfo() ([]util.StockInfo, error) {

	var stocks []util.StockInfo

	// 国内株式（現物/特定預り）
	sl, err := ss.GetStocksForJapanSpecificAccount()
	if err != nil {
		return nil, err
	}
	stocks = append(stocks, sl...)

	// 国内株式（現物/NISA預り）
	sl, err = ss.GetStocksForJapanNisaAccount()
	if err != nil {
		return nil, err
	}
	stocks = append(stocks, sl...)

	// 米国株式-全て
	sl, err = ss.GetStocksForUsaAccount()
	if err != nil {
		return nil, err
	}
	stocks = append(stocks, sl...)

	return stocks, nil
}

func (ss *SbiSecurities) openJapanAccountScreen() error {

	err := ss.ws.NavigatePage("https://site3.sbisec.co.jp/ETGate/?_ControlID=WPLETacR001Control&_PageID=DefaultPID&_ActionID=DefaultAID")
	if err != nil {
		return err
	}

	util.WaitTime()

	return nil
}

func (ss *SbiSecurities) GetStocksForJapanSpecificAccount() ([]util.StockInfo, error) {

	err := ss.openJapanAccountScreen()
	if err != nil {
		return nil, err
	}

	// 株式（現物/特定預り）
	multiSelection := ss.ws.GetPage().All("form table").At(1).All("table").At(17).All("tr")
	title, err := multiSelection.At(0).All("td > font > b").At(0).Text()
	if err != nil {
		return nil, err
	}
	if title != "株式（現物/特定預り）" {
		return nil, fmt.Errorf("構造が違います")
	}
	stocks := getStocksForJapan(multiSelection, util.SpecificAccount)
	return stocks, nil
}

func (ss *SbiSecurities) GetStocksForJapanNisaAccount() ([]util.StockInfo, error) {

	err := ss.openJapanAccountScreen()
	if err != nil {
		return nil, err
	}

	// 株式（現物/NISA預り）
	multiSelection := ss.ws.GetPage().All("form table").At(1).All("table").At(18).All("tr")
	title, err := multiSelection.At(0).All("td > font > b").At(0).Text()
	if err != nil {
		return nil, err
	}
	if title != "株式（現物/NISA預り）" {
		return nil, fmt.Errorf("構造が違います")
	}
	stocks := getStocksForJapan(multiSelection, util.NisaAccount)
	return stocks, nil
}

func getStocksForJapan(multiSelection *agouti.MultiSelection, securitiesAccount util.SecuritiesAccount) []util.StockInfo {

	count, _ := multiSelection.Count()
	count = (count - 2) / 2

	stockInfoList := make([]util.StockInfo, count)

	arrayCount := 0
	for i := 2; ; i++ {

		stockInfo := &util.StockInfo{
			SecuritiesCompany: util.SbiSecurities,
			StockCountry:      util.Japan,
			SecuritiesAccount: securitiesAccount,
		}

		// 奇数列は証券コードなど
		ms := multiSelection.At(i).All("td")

		// 証券コード
		secCode, err := ms.At(0).Text()
		if err != nil {
			break
		}
		stockInfo.SecuritiesCode = strconv.Itoa(util.ToIntByRemoveString(secCode))

		// 偶数列は保有株式数、取得単価など
		i++
		ms = multiSelection.At(i).All("td")

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

	return stockInfoList
}

func (ss *SbiSecurities) openUsaAccountScreen() error {

	err := ss.ws.NavigatePage("https://www.sbisec.co.jp/ETGate/?OutSide=on&_ControlID=WPLETsmR001Control&_DataStoreID=DSWPLETsmR001Control&sw_page=Foreign&cat1=home&cat2=none&sw_param1=GB&getFlg=on&int_pr1=150626_fstock_prodtop:odsite_btn_01")
	if err != nil {
		return err
	}

	util.WaitTime()

	err = ss.ws.NavigatePage("https://global.sbisec.co.jp/account/summary")
	if err != nil {
		return err
	}

	util.WaitTime()

	return nil
}

func (ss *SbiSecurities) GetStocksForUsaAccount() ([]util.StockInfo, error) {

	err := ss.openUsaAccountScreen()
	if err != nil {
		return nil, err
	}

	multiSelection := ss.ws.GetPage().All("div[class^=item-right] ul[class^=grid-table] > li")

	var stocks []util.StockInfo

	securitiesAccount := util.NoneAccount
	for i := 0; ; i++ {

		title, err := multiSelection.At(i).All("div > div").At(0).Text()
		if err != nil {
			break
		}
		if title == "株式(現物/特定)" {
			securitiesAccount = util.SpecificAccount
		} else if title == "株式(現物/一般)" {
			securitiesAccount = util.GeneralAccount
		} else if title == "株式(現物/NISA)" {
			securitiesAccount = util.NisaAccount
		}

		// 実際の株情報を見るために1つ進める
		i++
		sel := multiSelection.At(i)

		for j := 0; ; j++ {

			stockInfo := util.StockInfo{
				SecuritiesCompany: util.SbiSecurities,
				StockCountry:      util.Usa,
				SecuritiesAccount: securitiesAccount,
			}

			// 証券コード
			stockInfo.SecuritiesCode, err = sel.All("div > div[data-security-code]").At(j).Attribute("data-security-code")
			if err != nil {
				break
			}

			ms := sel.All("div > label")

			// 保有数量
			numOfStock, err := ms.At(j * 4).Text()
			if err != nil {
				break
			}
			stockInfo.NumberOfOwnedStock = util.ToIntByRemoveString(numOfStock)

			// 取得単価
			priceOfAvg, err := ms.At(j*4 + 1).Text()
			if err != nil {
				break
			}
			stockInfo.AveragePurchasePrice, err = util.ToFloatByRemoveString(priceOfAvg)
			if err != nil {
				break
			}

			stocks = append(stocks, stockInfo)
		}
	}

	return stocks, nil
}
