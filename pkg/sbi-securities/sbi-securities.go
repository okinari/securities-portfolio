package sbi_securities

import (
	"fmt"
	"github.com/okinari/golibs"
	"github.com/okinari/securities-portfolio/pkg/util"
	"github.com/sclevine/agouti"
	"strconv"
	"strings"
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

	url := "https://www.sbisec.co.jp/ETGate"
	err := ss.ws.NavigatePageForWait(url, 10*1000, 100)
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

	err = ss.ws.WaitForURLChange(url, 10*1000, 100)
	if err != nil {
		return err
	}

	return nil
}

func (ss *SbiSecurities) GetSecuritiesAccountInfo() ([]util.Stock, error) {

	var stocks []util.Stock

	// 国内株式（現物/特定預り）
	sl, err := ss.GetStocksForJapanSpecificAccount()
	if err != nil {
		return nil, err
	}
	stocks = append(stocks, sl...)

	// 国内株式（現物/NISA預り（成長投資枠））
	sl, err = ss.GetStocksForJapanNewNisaAccount()
	if err != nil {
		return nil, err
	}
	stocks = append(stocks, sl...)

	// 国内株式（現物/NISA預り（成長投資枠））
	sl, err = ss.GetStocksForJapanOldNisaAccount()
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
	err := ss.ws.NavigatePageForWait("https://site3.sbisec.co.jp/ETGate/?_ControlID=WPLETacR001Control&_PageID=DefaultPID&_ActionID=DefaultAID", 10*1000, 100)
	if err != nil {
		return err
	}

	return nil
}

func (ss *SbiSecurities) GetStocksForJapanSpecificAccount() ([]util.Stock, error) {

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

func (ss *SbiSecurities) GetStocksForJapanNewNisaAccount() ([]util.Stock, error) {

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
	if title != "株式（現物/NISA預り（成長投資枠））" {
		return nil, fmt.Errorf("構造が違います")
	}
	stocks := getStocksForJapan(multiSelection, util.NisaOldAccount)
	return stocks, nil
}

func (ss *SbiSecurities) GetStocksForJapanOldNisaAccount() ([]util.Stock, error) {

	err := ss.openJapanAccountScreen()
	if err != nil {
		return nil, err
	}

	// 株式（現物/NISA預り）
	multiSelection := ss.ws.GetPage().All("form table").At(1).All("table").At(19).All("tr")
	title, err := multiSelection.At(0).All("td > font > b").At(0).Text()
	if err != nil {
		return nil, err
	}
	if title != "株式（現物/旧NISA預り）" {
		return nil, fmt.Errorf("構造が違います")
	}
	stocks := getStocksForJapan(multiSelection, util.NisaOldAccount)
	return stocks, nil
}

func getStocksForJapan(multiSelection *agouti.MultiSelection, securitiesAccount util.SecuritiesAccount) []util.Stock {

	count, _ := multiSelection.Count()
	count = (count - 2) / 2

	stocks := make([]util.Stock, count)

	arrayCount := 0
	for i := 2; ; i++ {

		stock := &util.Stock{
			SecuritiesCompany: util.SbiSecurities,
			StockCountry:      util.Japan,
			SecuritiesAccount: securitiesAccount,
		}

		// 奇数列は証券コードなど
		ms := multiSelection.At(i).All("td")

		// 証券コード、会社名
		secCode, err := ms.At(0).Text()
		if err != nil {
			break
		}
		stock.SecuritiesCode = strconv.Itoa(util.ToIntByRemoveString(secCode))
		stock.CompanyName = strings.Replace(secCode, stock.SecuritiesCode, "", -1)
		stock.CompanyName = strings.TrimSpace(stock.CompanyName)

		// 偶数列は保有株式数、取得単価など
		i++
		ms = multiSelection.At(i).All("td")

		// 保有件数
		numOfStock, err := ms.At(0).Text()
		if err != nil {
			break
		}
		stock.NumberOfOwnedStock, err = util.ToFloatByRemoveString(numOfStock)
		if err != nil {
			break
		}

		// 取得単価
		priceOfAvg, err := ms.At(1).Text()
		if err != nil {
			break
		}
		stock.AveragePurchasePriceOne, err = util.ToFloatByRemoveString(priceOfAvg)
		if err != nil {
			break
		}

		stocks[arrayCount] = *stock
		arrayCount++
	}

	return stocks
}

func (ss *SbiSecurities) openUsaAccountScreen() error {

	url := "https://www.sbisec.co.jp/ETGate/?OutSide=on&_ControlID=WPLETsmR001Control&_DataStoreID=DSWPLETsmR001Control&sw_page=Foreign&cat1=home&cat2=none&sw_param1=GB&getFlg=on&int_pr1=150626_fstock_prodtop:odsite_btn_01"
	err := ss.ws.NavigatePage(url)
	if err != nil {
		return err
	}

	url = "https://global.sbisec.co.jp/home"
	err = ss.ws.WaitUntilURL(url, 10*1000, 1000)
	if err != nil {
		return err
	}

	url = "https://global.sbisec.co.jp/account/summary"
	err = ss.ws.NavigatePageForWait(url, 10*1000, 1000)
	if err != nil {
		return err
	}

	return nil
}

func (ss *SbiSecurities) GetStocksForUsaAccount() ([]util.Stock, error) {

	err := ss.openUsaAccountScreen()
	if err != nil {
		return nil, err
	}

	multiSelection := ss.ws.GetPage().All("div[class^=item-right] ul[class^=grid-table] > li")

	var stocks []util.Stock

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
			securitiesAccount = util.NisaOldAccount
		}

		// 実際の株情報を見るために1つ進める
		i++
		sel := multiSelection.At(i)

		for j := 0; ; j++ {

			stock := util.Stock{
				SecuritiesCompany: util.SbiSecurities,
				StockCountry:      util.Usa,
				SecuritiesAccount: securitiesAccount,
			}

			// 証券コード
			stock.SecuritiesCode, err = sel.All("div > div[data-security-code]").At(j).Attribute("data-security-code")
			if err != nil {
				break
			}

			ms := sel.All("div > label")

			// 保有数量
			numOfStock, err := ms.At(j * 4).Text()
			if err != nil {
				break
			}
			stock.NumberOfOwnedStock, err = util.ToFloatByRemoveString(numOfStock)
			if err != nil {
				break
			}

			// 取得単価
			priceOfAvg, err := ms.At(j*4 + 1).Text()
			if err != nil {
				break
			}
			stock.AveragePurchasePriceOne, err = util.ToFloatByRemoveString(priceOfAvg)
			if err != nil {
				break
			}

			stocks = append(stocks, stock)
		}
	}

	return stocks, nil
}