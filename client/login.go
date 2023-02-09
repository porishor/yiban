package client

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	url2 "net/url"
	"strings"
	"yiban/encrypt"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/thedevsaddam/gojsonq/v2"
)

// 构建易班客户端所需要的基本配置
func BuildRawYibanRequest(cookies []*http.Cookie) *resty.Client {
	req := resty.New().
		SetHeader("User-Agent", UserAgent).
		SetHeader("Origin", "https://c.uyiban.com").
		SetCookies(cookies)
	for _, c := range cookies {
		if c.Name == "csrf_token" {
			req.SetQueryParam("CSRF", c.Value)
		}
	}
	return req
}
func BuildYibanRequest(cookies []*http.Cookie) *resty.Request {
	return BuildRawYibanRequest(cookies).R()
}

// 初始化，返回cookie和跳转的网址
func initCookie() ([]*http.Cookie, *string, error) {
	var cookies = []*http.Cookie{
		{
			Name:  "csrf_token",
			Value: encrypt.GenCSRFToken(),
		},
	}
	resp, err := BuildYibanRequest(cookies).
		Get(ApiStartUp)
	if err != nil {
		return nil, nil, errors.New("error in constructing yiban client")
	} else {
		url := gojsonq.New().FromString(string(resp.Body())).Find("data.Data").(string)
		return append(resp.Cookies(), cookies...), &url, nil
	}
}

func getLoginPage(cookies *[]*http.Cookie, redirect_url string) (*resty.Response, error) {
	resp, err := BuildYibanRequest(*cookies).
		SetCookies(*cookies).
		Get(redirect_url)
	if err != nil {
		return nil, err
	} else {
		*cookies = append(*cookies, resp.Cookies()...)
		return resp, nil
	}
}

// 从cookies数组中提取出csrf_token
func ExtractCSRF(cookies []*http.Cookie) *string {
	return ExactCookie(cookies, "csrf_token")
}

func ExactCookie(cookies []*http.Cookie, key string) *string {
	var csrf_token = ""
	for _, cookie := range cookies {
		if cookie != nil {
			var k = cookie.Name
			if k == key {
				csrf_token = cookie.Value
				return &csrf_token
			}
		}
	}
	return nil
}

// 从网页中获得RSA公钥
func parseKeyFromPage(html []byte) (*string, error) {
	reader := strings.NewReader(string(html))
	dom, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	} else {
		key, exists := dom.Find("#key").Attr("value")
		if exists {
			return &key, nil
		} else {
			return nil, nil
		}
	}
}

// 验证身份,返回跳转的网址
func auth(redirect_uri string, cookies *[]*http.Cookie, username string, password []byte) (*string, error) {
	resp, _ := getLoginPage(cookies, redirect_uri)
	html := resp.Body()
	if html == nil {
		return nil, errors.New("Failed in getting login page")
	}
	key, _ := parseKeyFromPage(html)
	encrpted_password, err := encrypt.EncryptPassword(*key, []byte(password))
	if err != nil {
		return nil, errors.New("error in encryting password")
	}
	var pwd = base64.StdEncoding.EncodeToString(encrpted_password)
	resp2, err := BuildYibanRequest(*cookies).
		SetFormData(map[string]string{
			"oauth_uname":  username,
			"oauth_upwd":   pwd,
			"client_id":    "95626fa3080300ea",
			"redirect_uri": "https://f.yiban.cn/iapp7463",
			"state":        "",
			"scope":        "1,2,3,4,",
			"display":      "html",
		}).
		SetCookies(append(*cookies, resp.Cookies()...)).
		Post("https://oauth.yiban.cn/code/usersure")
	if err != nil {
		return nil, err
	}
	var res_text = string(resp2.Body())
	res_json := gojsonq.New().FromString(res_text)
	re_url := res_json.Find("reUrl").(string)
	// 返回结果不包含error，说明认证成功
	if !strings.Contains(re_url, "err") {
		*cookies = append(*cookies, resp2.Cookies()...)
		return &re_url, nil
	}
	return nil, errors.New(fmt.Sprintf("auth error:%s", resp2.Body()))
}

func Login(username string, password string) ([]*http.Cookie, error) {
	cookies, url, _ := initCookie()
	var csrf_token = ExtractCSRF(cookies)
	if csrf_token == nil {
		return nil, fmt.Errorf("Want csrf_token,but there is nothing in cookies\n")
	}
	url, err := auth(*url, &cookies, username, []byte(password))
	if err != nil {
		return nil, fmt.Errorf("验证身份失败\n")
	}
	var flag = false
	var verify_request = ""
	_, err = BuildRawYibanRequest(cookies).
		SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
			log.Printf("\n即将前往:%s\n", req.URL.String())
			if strings.Contains(req.URL.String(), "verify_request") {
				flag = true
				url, _ := url2.Parse(strings.ReplaceAll(req.URL.String(), "/#/", ""))
				verify_request = url.Query().Get("verify_request")
			}
			return nil
		})).
		R().
		Get(*url)
	BuildYibanRequest(cookies).
		SetQueryParam("verifyRequest", verify_request).
		Get(AuthURL)
	var new_cookies []*http.Cookie
	for _, c := range cookies {
		if c.Name == "csrf_token" || c.Name == "PHPSESSID" {
			new_cookies = append(new_cookies, c)
		}
	}
	if err != nil {
		return nil, err
	} else if !flag {
		log.Printf("\n登陆失败\n")
		return nil, err
	} else {
		log.Printf("登陆成功\n")
		return new_cookies, nil
	}
}

func BuildWithCookie(csrf_token string, PHPSESSID string) *resty.Request {
	return BuildYibanRequest([]*http.Cookie{
		{
			Name:  "csrf_token",
			Value: csrf_token,
		},
		{
			Name:  "PHPSESSID",
			Value: PHPSESSID,
		},
	})
}
