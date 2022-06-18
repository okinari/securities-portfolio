package sbi_securities

import (
	"fmt"
	"github.com/okinari/golibs"
	"github.com/okinari/securities-portfolio/pkg/util"
	"github.com/sclevine/agouti"
)

const LoginUrl = "https://www.sbisec.co.jp/ETGate"
const IpoUrl = "https://m.sbisec.co.jp/switchnaviMain"
const SecuritiesAccountUrl = "https://site3.sbisec.co.jp/ETGate/?_ControlID=WPLETacR001Control&_PageID=DefaultPID&_ActionID=DefaultAID"

type SbiSecurities struct {
	webScraping *golibs.WebScraping
}

func NewSbiSecurities() (*SbiSecurities, error) {
	ws, err := golibs.NewWebScraping(golibs.IsHeadless(false))
	if err != nil {
		return nil, err
	}

	return &SbiSecurities{
		webScraping: ws,
	}, nil
}

func (ss *SbiSecurities) Close() error {
	err := ss.webScraping.Close()
	if err != nil {
		return err
	}
	return nil
}

func (ss *SbiSecurities) Login(userName string, password string) error {

	err := ss.webScraping.NavigatePage(LoginUrl)
	if err != nil {
		return err
	}

	err = ss.webScraping.SetStringByName("user_id", userName)
	if err != nil {
		return err
	}
	err = ss.webScraping.SetStringByName("user_password", password)
	if err != nil {
		return err
	}
	err = ss.webScraping.ClickButtonByName("ACT_login")
	if err != nil {
		return err
	}

	return nil
}

func (ss *SbiSecurities) ApplyIpo() error {
	return nil
}

func (ss *SbiSecurities) GetSecuritiesAccountInfo() ([]util.StockInfo, error) {

	var stockInfoList []util.StockInfo

	err := ss.webScraping.NavigatePage(SecuritiesAccountUrl)
	if err != nil {
		return nil, err
	}

	page := ss.webScraping.GetPage()

	// 株式（現物/特定預り）
	multiSelection := page.All("form table").At(1).All("table").At(17).All("tr")
	title, err := multiSelection.At(0).All("td > font > b").At(0).Text()
	if err != nil {
		return nil, err
	}
	if title != "株式（現物/特定預り）" {
		return nil, fmt.Errorf("構造が違います")
	}
	siList := getStockList(multiSelection, util.SpecificAccount)
	stockInfoList = append(stockInfoList, siList...)

	// 株式（現物/NISA預り）
	multiSelection = page.All("form table").At(1).All("table").At(18).All("tr")
	title, err = multiSelection.At(0).All("td > font > b").At(0).Text()
	if err != nil {
		return nil, err
	}
	if title != "株式（現物/NISA預り）" {
		return nil, fmt.Errorf("構造が違います")
	}
	siList = getStockList(multiSelection, util.NisaAccount)
	stockInfoList = append(stockInfoList, siList...)

	//fmt.Printf("%v", stockInfoList)

	return stockInfoList, nil
}

func getStockList(multiSelection *agouti.MultiSelection, securitiesAccount util.SecuritiesAccount) []util.StockInfo {

	count, _ := multiSelection.Count()
	count = (count - 2) / 2

	stockInfoList := make([]util.StockInfo, count)
	//fmt.Printf("count: %v", count)

	arrayCount := 0
	for i := 2; ; i++ {

		stockInfo := &util.StockInfo{
			SecuritiesCompany: util.SbiSecurities,
			SecuritiesAccount: securitiesAccount,
		}

		// 奇数列は証券コードなど
		ms := multiSelection.At(i).All("td")

		// 証券コード
		secCode, err := ms.At(0).Text()
		if err != nil {
			//fmt.Printf("%v", err)
			break
		}
		stockInfo.SecuritiesCode = util.ToIntByRemoveString(secCode)
		//fmt.Printf("証券コード: %v", securitiesCode)

		// 偶数列は保有株式数、取得単価など
		i++
		ms = multiSelection.At(i).All("td")

		// 保有件数
		numOfStock, err := ms.At(0).Text()
		if err != nil {
			//fmt.Printf("%v", err)
			break
		}
		stockInfo.NumberOfOwnedStock = util.ToIntByRemoveString(numOfStock)
		//fmt.Printf("保有件数: %v", numOfStock)

		// 取得単価
		priceOfAvg, err := ms.At(1).Text()
		if err != nil {
			//fmt.Printf("%v", err)
			break
		}
		stockInfo.AveragePurchasePrice = util.ToIntByRemoveString(priceOfAvg)
		//fmt.Printf("取得単価: %v", priceOfAvg)

		stockInfoList[arrayCount] = *stockInfo
		arrayCount++
	}

	return stockInfoList
}
