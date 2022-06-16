package sbi_securities

import (
	"fmt"
	"github.com/okinari/golibs"
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

	//IpoUrl

	return nil
}

func (ss *SbiSecurities) GetSecuritiesAccountInfo() error {

	err := ss.webScraping.NavigatePage(SecuritiesAccountUrl)
	if err != nil {
		return err
	}

	page := ss.webScraping.GetPage()

	// 株式（現物/特定預り）
	multiSel := page.All("form table").At(1).All("table").At(17).All("tr")
	title, err := multiSel.At(0).All("td > font > b").At(0).Text()
	if err != nil {
		return err
	}

	if title != "株式（現物/特定預り）" {
		return fmt.Errorf("構造が違います")
	}

	for i := 2; ; i++ {
		ms := multiSel.At(i).All("td")
		for i := 0; ; i++ {
			str, err := ms.At(i).Text()
			if err != nil {
				break
			}

			fmt.Printf("result: %v \n", str)
		}

		_, err = ms.At(0).Text()
		if err != nil {
			break
		}

		fmt.Printf("ok \n")
	}

	// 株式（現物/NISA預り）
	multiSel = page.All("form table").At(1).All("table").At(18).All("tr")
	title, err = multiSel.At(0).All("td > font > b").At(0).Text()
	if err != nil {
		return err
	}

	if title != "株式（現物/NISA預り）" {
		return fmt.Errorf("構造が違います")
	}

	for i := 2; ; i++ {
		ms := multiSel.At(i).All("td")
		for i := 0; ; i++ {
			str, err := ms.At(i).Text()
			if err != nil {
				break
			}

			fmt.Printf("result: %v \n", str)
		}

		_, err = ms.At(0).Text()
		if err != nil {
			break
		}

		fmt.Printf("ok \n")
	}

	return nil
}
