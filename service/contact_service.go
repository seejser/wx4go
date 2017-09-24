package service

import (
	"encoding/json"
	"fmt"
	e "github.com/soease/wx4go/enum"
	m "github.com/soease/wx4go/model"
	t "github.com/soease/wx4go/tools"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

/**
 * 获取所有联系人信息，组装到map中，key为用户的UserName
 * 微信API对此URL使用了Cookie验证
 * GET:https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxgetcontact?lang=zh_CN&pass_ticket=dfLHy%252Fcgw%252BFM1qGhuARU6%252BDGs%252BGmWAD3jZJk6%252BfWcPs%253D&r=1504587952374&seq=0&skey=@crypt_3aaab8d5_c87634a7c5f8f579095cfdceeb8d842a
 */
func GetAllContact(loginMap *m.LoginMap) (map[string]m.User, error) {
	contactMap := map[string]m.User{}

	urlMap := e.GetInitParaEnum()
	urlMap[e.PassTicket] = loginMap.PassTicket
	urlMap[e.R] = fmt.Sprintf("%d", time.Now().UnixNano()/1000000)
	urlMap["seq"] = "0"
	urlMap[e.SKey] = loginMap.BaseRequest.SKey

	/* 使用Cookie功能，Get数据 */
	u, _ := url.Parse("https://wx.qq.com")

	jar := new(m.Jar)
	jar.SetCookies(u, loginMap.Cookies)

	client := &http.Client{
		Jar: jar}

	resp, err := client.Get(e.GET_ALL_CONTACT_URL + t.GetURLParams(urlMap))
	if err != nil {
		return contactMap, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return contactMap, err
	}

	contactList := m.ContactList{}
	err = json.Unmarshal(bodyBytes, &contactList)
	//fmt.Println(string(bodyBytes))
	if err != nil {
		return contactMap, err
	}

	for i := 0; i < contactList.MemberCount; i++ {
		contactMap[contactList.MemberList[i].UserName] = contactList.MemberList[i]
	}

	return contactMap, nil
}
