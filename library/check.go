package library

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	loginURL = url.URL{
		Scheme:   "http",
		Host:     "authserver.njc.ucas.ac.cn",
		Path:     "/authserver/login",
		RawQuery: "service=http%3A%2F%2Fehall.njc.ucas.ac.cn%2Fqljfwapp%2Fsys%2FlwPsXykApp%2Findex.do%23%2Fdledcx",
	}

	fetchURL = url.URL{
		Scheme: "http",
		Host:   "ehall.njc.ucas.ac.cn",
		Path:   "/qljfwapp/sys/lwPsXykApp/modules/dledcx/hqdled.do",
	}
)

type User struct {
	PushPlusClient `mapstructure:",squash"`
	Account        string `mapstructure:"account"`
	Password       string `mapstructure:"password"`
	RoomID         string `mapstructure:"room_id"`

	httpClient *http.Client
}

func (u *User) Init() {
	if u.httpClient == nil {
		jar, _ := cookiejar.New(nil)
		u.httpClient = &http.Client{
			Jar:     jar,
			Timeout: 15 * time.Second,
		}
	}
}

func (u *User) Check() {
	tryCounter := 1
	for i := range tryCounter {
		err := u.login()
		if err != nil {
			if errors.Is(err, needCaptchaErr) {
				logrus.WithField("user", u.RoomID).Warn("需要验证码，跳过本次查询")
				u.Send("查询前登录失败", "需要验证码，本次查询已跳过")
				return
			}
			if i == tryCounter-1 {
				logrus.WithField("user", u.RoomID).WithError(err).Error("登录失败, 达到最大重试次数")
				break
			}
			logrus.WithField("retried", i).WithField("user", u.RoomID).WithError(err).Warn("登录失败, 重试中...")
			continue
		}
		dataValue, err := u.fetchPowerData()
		if err != nil {
			if i == tryCounter-1 {
				logrus.WithField("user", u.RoomID).WithError(err).Error("查询电量失败, 达到最大重试次数")
				break
			}
			logrus.WithField("retried", i).WithField("user", u.RoomID).WithError(err).Warn("查询电量失败, 重试中...")
			continue
		}
		logrus.WithField("user", u.RoomID).WithField("current-power", dataValue).Info("查询电量成功")
		u.Send(fmt.Sprintf("当前电量: %s", dataValue), "电量查询成功")

		Insert(Metric{
			GID:   u.RoomID,
			Key:   "power",
			Value: dataValue,
		})

		return
	}

	u.Send("查询电量失败", "出现验证码拦截或其他异常")
}

var needCaptchaErr = errors.New("需要验证码，跳过本次查询")

func (u *User) captcha() (bool, error) {
	checkUrl := fmt.Sprintf("http://authserver.njc.ucas.ac.cn/authserver/checkNeedCaptcha.htl?username=%s&_=%d", u.Account, time.Now().UnixMilli())

	resp, err := u.httpClient.Get(checkUrl)
	if err != nil {
		return false, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var result struct {
		IsNeed bool `json:"isNeed"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, err
	}

	return result.IsNeed, nil
}

func (u *User) login() error {
	u_url, _ := url.Parse(loginURL.String())
	u.httpClient.Jar.SetCookies(u_url, []*http.Cookie{{
		Name:  "org.springframework.web.servlet.i18n.CookieLocaleResolver.LOCALE",
		Value: "zh_CN",
		Path:  "/",
	}})

	resp, err := u.httpClient.Get(loginURL.String())
	if err != nil {
		return err
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	html := string(body)

	// 检查是否需要验证码
	needCaptcha, err := u.captcha()
	if err != nil {
		return err
	}
	if needCaptcha {
		return needCaptchaErr
	}

	// 提取 execution 和 encryptPassword salt
	execRe := regexp.MustCompile(`name="execution" value="([^"]+)"`)
	execMatch := execRe.FindStringSubmatch(html)
	if len(execMatch) < 2 {
		return errors.New("未能提取 execution")
	}
	executation := execMatch[1]
	saltRe := regexp.MustCompile(`id="pwdEncryptSalt" value="([^"]+)"`)
	saltMatch := saltRe.FindStringSubmatch(html)
	if len(saltMatch) < 2 {
		return errors.New("未能提取 pwdEncryptSalt")
	}
	salt := saltMatch[1]

	// 加密密码
	encPassword, err := encryptPassword(u.Password, salt)
	if err != nil {
		return err
	}

	form := url.Values{}
	form.Set("username", u.Account)
	form.Set("password", encPassword)
	form.Set("lt", "")
	form.Set("dllt", "generalLogin")
	form.Set("execution", executation)
	form.Set("_eventId", "submit")
	form.Set("cllt", "userNameLogin")

	req, err := http.NewRequest("POST", loginURL.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:146.0) Gecko/20100101 Firefox/146.0")
	req.Header.Set("Referer", loginURL.String())

	postResp, err := u.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer postResp.Body.Close()

	logrus.WithField("cookie", u.httpClient.Jar.Cookies(&loginURL)).Debug("登录成功")
	return nil
}

// fetchPowerData 请求获取电量信息
func (u *User) fetchPowerData() (string, error) {
	formData := url.Values{}

	// 创建请求，将表单数据作为 Body
	req, err := http.NewRequest("POST", fetchURL.String(), strings.NewReader(formData.Encode()))
	if err != nil {
		return "", err
	}

	// 3. 添加必要的请求头
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // 必须指定
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:146.0) Gecko/20100101 Firefox/146.0")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Referer", "http://ehall.njc.ucas.ac.cn/qljfwapp/sys/lwPsXykApp/index.do")

	// 4. 发送请求
	resp, err := u.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("服务器响应异常: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 5. 解析响应 (保持你原有的结构)
	var result struct {
		Datas struct {
			Hqdled struct {
				Rows []struct {
					REMAINEQ float64 `json:"REMAINEQ"`
				} `json:"rows"`
			} `json:"hqdled"`
		} `json:"datas"`
		Code string `json:"code"` // 建议增加 code 字段判断
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	// 校验业务逻辑
	if len(result.Datas.Hqdled.Rows) == 0 {
		return "", errors.New("未找到电量信息，请检查登录状态")
	}

	dataValue := fmt.Sprintf("%.2f", result.Datas.Hqdled.Rows[0].REMAINEQ)
	return dataValue, nil
}
