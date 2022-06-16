package sbi_neomobile

import (
	"fmt"
	"github.com/okinari/golibs"
	"github.com/okinari/securities-portfolio/pkg/util"
	"strconv"
	"strings"
)

const LoginUrl = "https://trade.sbineomobile.co.jp/login"
const SecuritiesAccountUrl = "https://trade.sbineomobile.co.jp/account/portfolio"

type SbiNeomobile struct {
	webScraping *golibs.WebScraping
}

func NewSbiNeomobile() (*SbiNeomobile, error) {
	ws, err := golibs.NewWebScraping(golibs.IsHeadless(false))
	if err != nil {
		return nil, err
	}

	return &SbiNeomobile{
		webScraping: ws,
	}, nil
}

func (sn *SbiNeomobile) Close() error {
	err := sn.webScraping.Close()
	if err != nil {
		return err
	}
	return nil
}

func (sn *SbiNeomobile) Login(userName string, password string) error {

	err := sn.webScraping.NavigatePage(LoginUrl)
	if err != nil {
		return err
	}

	err = sn.webScraping.SetStringByName("username", userName)
	if err != nil {
		return err
	}
	err = sn.webScraping.SetStringByClass("input-password", password)
	if err != nil {
		return err
	}
	err = sn.webScraping.ClickButtonByID("neo-login-btn")
	if err != nil {
		return err
	}

	return nil
}

func (sn *SbiNeomobile) ApplyIpo() error {
	return nil
}

func (sn *SbiNeomobile) GetSecuritiesAccountInfo() ([]util.StockInfo, error) {

	err := sn.webScraping.NavigatePage(SecuritiesAccountUrl)
	if err != nil {
		return nil, err
	}
	util.WaitTime()

	page := sn.webScraping.GetPage()

	// 最後まで表示
	for {
		err = sn.webScraping.ClickButtonByClass("more")
		if err != nil {
			break
		}
	}

	// 証券番号
	// document.querySelectorAll('div.panel')[n].querySelector('.ticker')

	// 詳細の項目名
	// document.querySelectorAll('div.panel')[n].querySelectorAll('table tbody tr')[n].querySelector('th')

	// 詳細の値
	// document.querySelectorAll('div.panel')[n].querySelectorAll('table tbody tr')[n].querySelector('td span')

	// 株式（現物/特定預り）
	multiSel := page.All("div.panel")

	// 全部開く
	for i := 0; ; i++ {
		err = multiSel.At(i).Find("div.price").Click()
		if err != nil {
			break
		}
	}

	var stockInfoList []util.StockInfo

	for i := 0; ; i++ {

		stock := util.StockInfo{
			SecuritiesCompany: util.SbiNeomobile,
		}

		// 証券コード
		secCode, err := multiSel.At(i).FindByClass("ticker").Text()
		if err != nil {
			//fmt.Printf("error: 証券コード: %v \n", err)
			break
		}
		securitiesCode, err := strconv.Atoi(secCode)
		if err != nil {
			fmt.Printf("error: 数値変換に失敗: %v \n", err)
			break
		}
		stock.SecuritiesCode = securitiesCode
		//fmt.Printf("result 証券番号: %v \n", securitiesCode)

		ms := multiSel.At(i).All("table tbody tr")

		for j := 0; ; j++ {
			item := ms.At(j)
			//fmt.Printf("result: 項目: %v \n", item)

			name, err := item.Find("th").Text()
			if err != nil {
				//fmt.Printf("error: 項目名: %v \n", err)
				break
			}
			//fmt.Printf("result 項目名: %v \n", name)

			value, err := item.All("td > span").At(0).Text()
			if err != nil {
				//fmt.Printf("error: 項目値: %v \n", err)
				break
			}
			//fmt.Printf("result 項目値: %v \n", value)

			if name == "保有数量" {
				numberOfOwnedStock, err := strconv.Atoi(value)
				if err != nil {
					fmt.Printf("error: 数値変換に失敗: %v \n", err)
					break
				}
				stock.NumberOfOwnedStock = numberOfOwnedStock
				//fmt.Printf("result 保有数量: %v \n", value)
			} else if name == "平均取得単価" {
				value = strings.Replace(value, ",", "", -1)
				averagePurchasePrice, err := strconv.Atoi(value)
				if err != nil {
					fmt.Printf("error: 数値変換に失敗: %v \n", err)
					break
				}
				stock.AveragePurchasePrice = averagePurchasePrice
				//fmt.Printf("result 平均取得単価: %v \n", value)
			}
		}

		stockInfoList = append(stockInfoList, stock)
	}

	return stockInfoList, nil
}
