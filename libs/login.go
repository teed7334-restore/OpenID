package libs

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/RangelReale/osin"
)

//Login 相關設定
type Login struct {
	user     string
	password string
}

//New 建構式
func (l Login) New() *Login {
	return &l
}

//HandleLoginPage 登入頁面
func (l *Login) HandleLoginPage(ar *osin.AuthorizeRequest, w http.ResponseWriter, r *http.Request) *Login {
	r.ParseForm()
	l.user = ""
	l.password = ""

	if r.Method == "POST" {
		l.user = r.FormValue("user")
		l.password = r.FormValue("password")
		return l
	}

	w.Write([]byte("<html><body>"))

	w.Write([]byte(fmt.Sprintf("LOGIN %s (use test/test)<br/>", ar.Client.GetId())))
	w.Write([]byte(fmt.Sprintf("<form action=\"/authorize?%s\" method=\"POST\">", r.URL.RawQuery)))

	w.Write([]byte("Login: <input type=\"text\" name=\"user\" /><br/>"))
	w.Write([]byte("Password: <input type=\"password\" name=\"password\" /><br/>"))
	w.Write([]byte("<input type=\"submit\"/>"))

	w.Write([]byte("</form>"))

	w.Write([]byte("</body></html>"))

	return l
}

//DownloadAccessToken 下載存取Token頁面
func (l *Login) DownloadAccessToken(url string, auth *osin.BasicAuth, output map[string]interface{}) error {
	// download access token
	preq, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	if auth != nil {
		preq.SetBasicAuth(auth.Username, auth.Password)
	}

	pclient := &http.Client{}
	presp, err := pclient.Do(preq)
	if err != nil {
		return err
	}

	if presp.StatusCode != 200 {
		return errors.New("Invalid status code")
	}

	jdec := json.NewDecoder(presp.Body)
	err = jdec.Decode(&output)
	return err
}
