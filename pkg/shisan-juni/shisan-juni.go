package shisan_juni

import (
	"fmt"
	"github.com/okinari/golibs"
	"github.com/okinari/securities-portfolio/pkg/util"
	"strconv"
)

const LoginUrl = "https://43juni.pocco.net/auth/login"
const JapanSecuritiesUrl = "https://43juni.pocco.net/app/dashboard/japan"
const UsaSecuritiesUrl = "https://43juni.pocco.net/app/dashboard/usa"

type ShisanJuni struct {
	ws *golibs.WebScraping
}

func NewShisanJuni() (*ShisanJuni, error) {
	ws, err := golibs.NewWebScraping(golibs.IsHeadless(false))
	if err != nil {
		return nil, err
	}

	return &ShisanJuni{
		ws: ws,
	}, nil
}

func (sj *ShisanJuni) Close() error {
	err := sj.ws.Close()
	if err != nil {
		return err
	}
	return nil
}

func (sj *ShisanJuni) Login(userName string, password string) error {

	err := sj.ws.NavigatePage(LoginUrl)
	if err != nil {
		return err
	}

	util.WaitTime()

	err = sj.ws.SetStringByName("email", userName)
	if err != nil {
		return err
	}
	err = sj.ws.SetStringByName("password", password)
	if err != nil {
		return err
	}
	err = sj.ws.GetPage().Find("button[type=submit]").Click()
	if err != nil {
		return err
	}

	util.WaitTime()

	return nil
}

func (sj *ShisanJuni) ApplyIpo() error {

	//IpoUrl

	return nil
}

func (sj *ShisanJuni) SetSecuritiesInfo() error {

	err := sj.ws.NavigatePage(JapanSecuritiesUrl)
	if err != nil {
		return err
	}
	util.WaitTime()

	page := sj.ws.GetPage()

	// 次のページへ
	// document.querySelectorAll("[title='Go to next page']")[0].click()

	// 次のページの存在確認
	// document.querySelectorAll("[title='Go to next page']")[0].disabled

	// ループ回す対象
	// document.querySelectorAll("table")[0].querySelectorAll("tbody > tr")

	// 証券コード
	// document.querySelectorAll("table")[0].querySelectorAll("tbody > tr")[n].querySelectorAll("td")[1].querySelectorAll("div")[0].querySelector("span")

	// 証券会社
	// document.querySelectorAll("table")[0].querySelectorAll("tbody > tr")[0].querySelectorAll("td")[2].querySelectorAll('span span')[0]

	// 保有数
	// document.querySelectorAll("table")[0].querySelectorAll("tbody > tr")[n].querySelectorAll("td")[4]

	// 平均購入単価
	// document.querySelectorAll("table")[0].querySelectorAll("tbody > tr")[n].querySelectorAll("td")[5]

	multiSel := page.Find("table").All("tbody > tr")

	for i := 0; ; i++ {

		ms := multiSel.At(i).All("td")

		// 謎の行がある
		count, err := ms.Count()
		// カウントできない場合、そのページは終わり
		if err != nil {
			fmt.Printf("error: カウントできない: %v \n", err)

			//time.Sleep(SleepTime * time.Second)

			btn := page.All("button[title='Go to next page']").At(0)
			fmt.Printf("result: ボタン: %v \n", btn)

			// 次のページを確認し、存在する場合、次のページへ、存在しない場合、処理終了
			disabled, err := btn.Attribute("disabled")
			if err != nil {
				fmt.Printf("error: 次のページの確認に失敗: %v \n", err)
			}
			fmt.Printf("result: 次のページの存在確認: %v \n", disabled)
			if disabled == "true" {
				break
			}

			// 存在する場合
			err = sj.ws.ExecJavaScript("document.querySelectorAll(\"[title='Go to next page']\")[0].click()", nil)
			if err != nil {
				fmt.Printf("error: 次のページへの遷移に失敗: %v \n", err)
				break
			}

			i = -1

			continue
		}
		if count == 1 {
			continue
		}

		// 証券コード
		securitiesCode, err := ms.At(1).All("div").At(0).Find("span").Text()
		if err != nil {
			fmt.Printf("error: 証券コード: %v \n", err)
			break
		}
		fmt.Printf("result 証券コード: %v \n", securitiesCode)

		// 証券会社
		securitiesCompany, err := ms.At(2).All("span > span").At(0).Text()
		if err != nil {
			fmt.Printf("error: 証券会社: %v \n", err)
			break
		}
		fmt.Printf("result 証券会社: %v \n", securitiesCompany)

		// 保有株式数
		numOfOwnedStock, err := ms.At(4).Text()
		if err != nil {
			fmt.Printf("error: 保有株式数: %v \n", err)
			break
		}
		fmt.Printf("result 保有株式数: %v \n", numOfOwnedStock)

		// 平均購入単価
		unitPriceOfAvgPurchase, err := ms.At(5).Text()
		if err != nil {
			fmt.Printf("error: 平均購入単価: %v \n", err)
			break
		}
		fmt.Printf("result 平均購入単価: %v \n", unitPriceOfAvgPurchase)

	}

	return nil
}

func (sj *ShisanJuni) addSec(stockInfo util.StockInfo) error {

	err := sj.ws.NavigatePage(JapanSecuritiesUrl)
	if err != nil {
		return err
	}
	util.WaitTime()

	// 登録フロートを開く
	// document.querySelectorAll("div div div button[class^='MuiButtonBase-root MuiFab-root MuiFab-circular MuiFab-sizeLarge MuiFab-primary'] span[class^=MuiFab-label]")[0].click()

	// 証券会社
	// name=sCompany

	// 証券コード
	// name=stockCode

	// 追加株式数
	// name=addNumber

	// 約定単価
	// name=buyPrice

	err = sj.ws.ClickButtonBySelector("div > div > div > button[class^='MuiButtonBase-root MuiFab-root MuiFab-circular MuiFab-sizeLarge MuiFab-primary'] > span[class^=MuiFab-label]")
	if err != nil {
		return err
	}

	//sj.ws.SetStringByName("sCompany", stockInfo.SecuritiesCompany)
	err = sj.ws.SetStringByName("sCompany", "SBI証券")
	if err != nil {
		return err
	}
	err = sj.ws.SetStringByName("stockCode", strconv.Itoa(stockInfo.SecuritiesCode))
	if err != nil {
		return err
	}
	err = sj.ws.SetStringByName("addNumber", strconv.Itoa(stockInfo.NumberOfOwnedStock))
	if err != nil {
		return err
	}
	err = sj.ws.SetStringByName("buyPrice", strconv.Itoa(stockInfo.AveragePurchasePrice))
	if err != nil {
		return err
	}

	return nil
}

func (sj *ShisanJuni) UpdateSec(stockInfo util.StockInfo) error {

	err := sj.ws.NavigatePage(JapanSecuritiesUrl)
	if err != nil {
		return err
	}
	util.WaitTime()

	page := sj.ws.GetPage()

	// 次のページへ
	// document.querySelectorAll("[title='Go to next page']")[0].click()

	// 次のページの存在確認
	// document.querySelectorAll("[title='Go to next page']")[0].disabled

	// ループ回す対象
	// document.querySelectorAll("table")[0].querySelectorAll("tbody > tr")

	// 証券コード
	// document.querySelectorAll("table")[0].querySelectorAll("tbody > tr")[n].querySelectorAll("td")[1].querySelectorAll("div")[0].querySelector("span")

	// 証券会社
	// document.querySelectorAll("table")[0].querySelectorAll("tbody > tr")[0].querySelectorAll("td")[2].querySelectorAll('span span')[0]

	// 保有数
	// document.querySelectorAll("table")[0].querySelectorAll("tbody > tr")[n].querySelectorAll("td")[4]

	// 平均購入単価
	// document.querySelectorAll("table")[0].querySelectorAll("tbody > tr")[n].querySelectorAll("td")[5]

	multiSel := page.Find("table").All("tbody > tr")

	for i := 0; ; i++ {

		ms := multiSel.At(i).All("td")

		// 謎の行がある
		count, err := ms.Count()
		// カウントできない場合、そのページは終わり
		if err != nil {
			fmt.Printf("error: カウントできない: %v \n", err)

			//time.Sleep(SleepTime * time.Second)

			btn := page.All("button[title='Go to next page']").At(0)
			fmt.Printf("result: ボタン: %v \n", btn)

			// 次のページを確認し、存在する場合、次のページへ、存在しない場合、処理終了
			disabled, err := btn.Attribute("disabled")
			if err != nil {
				fmt.Printf("error: 次のページの確認に失敗: %v \n", err)
			}
			fmt.Printf("result: 次のページの存在確認: %v \n", disabled)
			if disabled == "true" {
				break
			}

			// 存在する場合
			err = sj.ws.ExecJavaScript("document.querySelectorAll(\"[title='Go to next page']\")[0].click()", nil)
			if err != nil {
				fmt.Printf("error: 次のページへの遷移に失敗: %v \n", err)
				break
			}

			i = -1

			continue
		}
		if count == 1 {
			continue
		}

		// 証券コード
		secCode, err := ms.At(1).All("div").At(0).Find("span").Text()
		if err != nil {
			fmt.Printf("error: 証券コード: %v \n", err)
			break
		}
		//fmt.Printf("result 証券コード: %v \n", securitiesCode)
		securitiesCode, err := strconv.Atoi(secCode)
		if err != nil {
			fmt.Printf("error: 証券コードの数字化に失敗: %v \n", err)
			break
		}

		// 証券会社
		secCompany, err := ms.At(2).All("span > span").At(0).Text()
		if err != nil {
			fmt.Printf("error: 証券会社: %v \n", err)
			break
		}
		//fmt.Printf("result 証券会社: %v \n", securitiesCompany)

		// 該当の株式ではない場合、次へ
		securitiesCompany := util.GetSecuritiesCompany(secCompany)
		if stockInfo.SecuritiesCode != securitiesCode || stockInfo.SecuritiesCompany != securitiesCompany {
			continue
		}

		fmt.Printf("該当株式発見")
		util.WaitTime()

		// 情報更新のためのポップアップを開く
		err = sj.ws.ExecJavaScript("document.querySelector('table > tbody').scrollIntoView()", nil)
		if err != nil {
			fmt.Printf("error: スクロール処理に失敗: %v \n", err)
			break
		}
		err = ms.At(7).All("div > button > span").At(4).Click()
		if err != nil {
			fmt.Printf("error: 情報更新のためのポップアップを開く操作に失敗: %v \n", err)
			break
		}

		fmt.Printf("ポップアップ開く")
		util.WaitTime()

		// 保有株式数、約定単価を設定して、情報更新
		err = sj.ws.ExecJavaScript("document.querySelector('[name=editPosAmount]').value = '"+strconv.Itoa(stockInfo.NumberOfOwnedStock)+"'", nil)
		if err != nil {
			fmt.Printf("error: 保有株式数のクリアに失敗: %v \n", err)
			break
		}
		err = sj.ws.SetStringByName("editPosAmount", "")
		err = sj.ws.SetStringByName("editPosAmount", strconv.Itoa(stockInfo.NumberOfOwnedStock))
		if err != nil {
			fmt.Printf("error: 保有株式数の設定に失敗: %v \n", err)
			break
		}
		err = sj.ws.ExecJavaScript("document.querySelector('[name=editAvgPrice]').value = '"+strconv.Itoa(stockInfo.AveragePurchasePrice)+"'", nil)
		if err != nil {
			fmt.Printf("error: 約定単価のクリアに失敗: %v \n", err)
			break
		}
		err = sj.ws.SetStringByName("editAvgPrice", "")
		err = sj.ws.SetStringByName("editAvgPrice", strconv.Itoa(stockInfo.AveragePurchasePrice))
		if err != nil {
			fmt.Printf("error: 約定単価の設定に失敗: %v \n", err)
			break
		}

		util.WaitTime()

		err = sj.ws.ClickButtonBySelector("form button[type=submit]")
		if err != nil {
			fmt.Printf("error: 更新処理に失敗: %v \n", err)
			break
		}

	}

	// 次のページへ
	// document.querySelectorAll("[title='Go to next page']")[0].click()

	// ループ回す対象
	// document.querySelectorAll("table")[0].querySelectorAll("tbody > tr")

	// 証券コード
	// document.querySelectorAll("table")[0].querySelectorAll("tbody > tr")[n].querySelectorAll("td")[1].querySelectorAll("div")[0].querySelector("span")

	// 証券会社
	// document.querySelectorAll("table")[0].querySelectorAll("tbody > tr")[0].querySelectorAll("td")[2].querySelectorAll('span span')[0]

	// 更新フロートを開く
	// document.querySelectorAll("table")[0].querySelectorAll("tbody > tr")[n].querySelectorAll("td")[7].querySelectorAll('div button')[2].click()

	// 保有株式数
	// name=editPosAmount

	// 約定単価
	// name=editAvgPrice

	// 編集ボタンを押下
	// document.querySelectorAll('form button[type=submit]')[0].click()

	return nil
}
