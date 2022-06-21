package sbi_neomobile

import (
	"github.com/okinari/golibs"
	"github.com/okinari/securities-portfolio/pkg/util"
	"strconv"
)

type SbiNeomobile struct {
	ws *golibs.WebScraping
}

func NewSbiNeomobile() (*SbiNeomobile, error) {
	ws, err := golibs.NewWebScraping(golibs.IsHeadless(false))
	if err != nil {
		return nil, err
	}

	return &SbiNeomobile{
		ws: ws,
	}, nil
}

func (sn *SbiNeomobile) Close() error {
	err := sn.ws.Close()
	if err != nil {
		return err
	}
	return nil
}

func (sn *SbiNeomobile) Login(userName string, password string) error {

	err := sn.ws.NavigatePage("https://trade.sbineomobile.co.jp/login")
	if err != nil {
		return err
	}

	err = sn.ws.SetStringByName("username", userName)
	if err != nil {
		return err
	}
	err = sn.ws.SetStringByClass("input-password", password)
	if err != nil {
		return err
	}
	err = sn.ws.ClickButtonByID("neo-login-btn")
	if err != nil {
		return err
	}

	return nil
}

func (sn *SbiNeomobile) GetSecuritiesAccountInfo() ([]util.StockInfo, error) {

	err := sn.ws.NavigatePage("https://trade.sbineomobile.co.jp/account/portfolio")
	if err != nil {
		return nil, err
	}
	util.WaitTime()

	// 最後まで表示
	for {
		err = sn.ws.ClickButtonByClass("more")
		if err != nil {
			break
		}
	}

	// 株式（現物/特定預り）
	multiSelection := sn.ws.GetPage().All("div.panel")

	// 全部開く
	for i := 0; ; i++ {
		err = multiSelection.At(i).Find("div.price").Click()
		if err != nil {
			break
		}
	}

	var stockInfoList []util.StockInfo

	for i := 0; ; i++ {

		stock := util.StockInfo{
			SecuritiesCompany: util.SbiNeomobile,
			StockCountry:      util.Japan,
			SecuritiesAccount: util.SpecificAccount,
		}

		// 証券コード
		secCode, err := multiSelection.At(i).FindByClass("ticker").Text()
		if err != nil {
			break
		}
		stock.SecuritiesCode = secCode

		ms := multiSelection.At(i).All("table tbody tr")

		for j := 0; ; j++ {
			item := ms.At(j)

			name, err := item.Find("th").Text()
			if err != nil {
				break
			}

			value, err := item.All("td > span").At(0).Text()
			if err != nil {
				break
			}

			if name == "保有数量" {
				numberOfOwnedStock, err := strconv.Atoi(value)
				if err != nil {
					break
				}
				stock.NumberOfOwnedStock = numberOfOwnedStock
			} else if name == "平均取得単価" {
				stock.AveragePurchasePrice, err = util.ToFloatByRemoveString(value)
				if err != nil {
					break
				}
			}
		}

		stockInfoList = append(stockInfoList, stock)
	}

	return stockInfoList, nil
}
