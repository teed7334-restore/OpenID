package libs

import (
	"net/http"

	"github.com/RangelReale/osin"
)

//OAuth 相關設定
type OAuth struct {
	server *osin.Server
}

//New 建構式
func (oh OAuth) New() OAuth {
	storage := Storage{}.New()
	oh.server = osin.NewServer(osin.NewServerConfig(), storage)
	return oh
}

//APIs 呼叫的API路徑列表
func (oh *OAuth) APIs() {
	http.HandleFunc("/authorize", oh.authorize)
	http.HandleFunc("/token", oh.token)
	http.ListenAndServe(":14000", nil)
}

func (oh *OAuth) open(r *http.Request) (*osin.Response, *osin.AuthorizeRequest) {
	resp := oh.server.NewResponse()
	defer resp.Close()
	ar := oh.server.HandleAuthorizeRequest(resp, r)
	return resp, ar
}

func (oh *OAuth) authorize(w http.ResponseWriter, r *http.Request) {
	resp, ar := oh.open(r)

	if ar != nil {
		l := Login{}.New()
		result := l.HandleLoginPage(ar, w, r)
		if (result.user == "" || result.password == "") || !oh.authenticate(result.user, result.password) {
			return
		}
		ar.Authorized = true
		oh.server.FinishAuthorizeRequest(resp, r, ar)
	}
	osin.OutputJSON(resp, w, r)
}

func (oh *OAuth) authenticate(user, password string) bool {
	l := LDAP{}.New()
	conn := l.Connect()
	defer l.Close(conn)
	validated := l.Login(conn, user, password)
	return validated
}

func (oh *OAuth) token(w http.ResponseWriter, r *http.Request) {
	resp, ar := oh.open(r)
	if ar != nil {
		ar.Authorized = true
		oh.server.FinishAuthorizeRequest(resp, r, ar)
	}
	osin.OutputJSON(resp, w, r)
}
