package libs

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/ldap.v3"
)

//LDAP 相關設定
type LDAP struct {
	host     string
	port     string
	user     string
	password string
	dc       string
}

//New 建構式
func (lp LDAP) New() *LDAP {
	lp.host = os.Getenv("LDAP_HOST_NAME")
	lp.port = os.Getenv("LDAP_HOST_PORT")
	lp.user = os.Getenv("LDAP_USER_NAME")
	lp.password = os.Getenv("LDAP_USER_PASSWORD")
	lp.dc = os.Getenv("LDAP_DC_NAME")
	return &lp
}

//Connect 連線到LDAP
func (lp *LDAP) Connect() *ldap.Conn {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%s", lp.host, lp.port))
	if err != nil {
		log.Panicln(err)
	}

	err = l.Bind(fmt.Sprintf("cn=%s,%s", lp.user, lp.dc), lp.password)
	if err != nil {
		log.Panicln(err)
	}
	return l
}

//Close 關閉LDAP連線
func (lp *LDAP) Close(l *ldap.Conn) {
	l.Close()
}

//Login 透過使用者登入LDAP
func (lp *LDAP) Login(l *ldap.Conn, user string, password string) bool {
	searchRequest := ldap.NewSearchRequest(
		lp.dc,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", user),
		[]string{"dn"},
		nil,
	)
	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Panicln(err)
	}
	userdn := sr.Entries[0].DN
	err = l.Bind(userdn, password)
	if err != nil {
		log.Panicln(err)
		return false
	}
	return true
}
