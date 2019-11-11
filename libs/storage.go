package libs

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/RangelReale/osin"
)

//oauthClients 人資系統OAuth Client資料表結構
type oauthClients struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	GrantTypes   string
	Scope        string
	UserID       string
}

//Storage OAuth 資料結構
type Storage struct {
	clients   map[string]osin.Client
	authorize map[string]*osin.AuthorizeData
	access    map[string]*osin.AccessData
	refresh   map[string]string
}

//New 建構式
func (s Storage) New() *Storage {
	s.clients = make(map[string]osin.Client)
	clients := s.getClients()
	expire64, err := strconv.Atoi(os.Getenv("EXPIRE_TIME"))
	expire := int32(expire64)
	if err != nil {
		log.Panicln(err)
	}
	for _, item := range clients {
		s.clients[item.ClientID] = &osin.DefaultClient{
			Id:          item.ClientID,
			Secret:      item.ClientSecret,
			RedirectUri: item.RedirectURI,
		}
		s.authorize = make(map[string]*osin.AuthorizeData)
		s.access = make(map[string]*osin.AccessData)
		s.refresh = make(map[string]string)
		s.access[item.ClientID] = &osin.AccessData{
			Client:        s.clients[item.ClientID],
			AuthorizeData: s.authorize[item.ClientID],
			AccessToken:   item.ClientID,
			ExpiresIn:     expire,
			CreatedAt:     time.Now(),
		}
	}
	return &s
}

//Clone 複製新的Client
func (s *Storage) Clone() osin.Storage {
	return s
}

//Close 關閉Client
func (s *Storage) Close() {
}

//GetClient 取得Client
func (s *Storage) GetClient(id string) (osin.Client, error) {
	fmt.Printf("GetClient: %s\n", id)
	if c, ok := s.clients[id]; ok {
		return c, nil
	}
	return nil, osin.ErrNotFound
}

//SetClient 設定Client
func (s *Storage) SetClient(id string, client osin.Client) error {
	fmt.Printf("SetClient: %s\n", id)
	s.clients[id] = client
	return nil
}

//SaveAuthorize 儲存授權許可
func (s *Storage) SaveAuthorize(data *osin.AuthorizeData) error {
	fmt.Printf("SaveAuthorize: %s\n", data.Code)
	s.saveAuthorize2Redis(data.Code, "1")
	s.authorize[data.Code] = data
	return nil
}

//LoadAuthorize 讀取授權許可
func (s *Storage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	fmt.Printf("LoadAuthorize: %s\n", code)
	if d, ok := s.authorize[code]; ok {
		return d, nil
	}
	return nil, osin.ErrNotFound
}

//RemoveAuthorize 移除授權許可
func (s *Storage) RemoveAuthorize(code string) error {
	fmt.Printf("RemoveAuthorize: %s\n", code)
	delete(s.authorize, code)
	return nil
}

//SaveAccess 儲存存取資訊
func (s *Storage) SaveAccess(data *osin.AccessData) error {
	fmt.Printf("SaveAccess: %s\n", data.AccessToken)
	s.access[data.AccessToken] = data
	if data.RefreshToken != "" {
		s.refresh[data.RefreshToken] = data.AccessToken
	}
	return nil
}

//LoadAccess 載入存取資訊
func (s *Storage) LoadAccess(code string) (*osin.AccessData, error) {
	fmt.Printf("LoadAccess: %s\n", code)
	if d, ok := s.access[code]; ok {
		return d, nil
	}
	return nil, osin.ErrNotFound
}

//RemoveAccess 移除存取資訊
func (s *Storage) RemoveAccess(code string) error {
	fmt.Printf("RemoveAccess: %s\n", code)
	delete(s.access, code)
	return nil
}

//LoadRefresh 讀入換發Token
func (s *Storage) LoadRefresh(code string) (*osin.AccessData, error) {
	fmt.Printf("LoadRefresh: %s\n", code)
	if d, ok := s.refresh[code]; ok {
		return s.LoadAccess(d)
	}
	return nil, osin.ErrNotFound
}

//RemoveRefresh 移除換發Token
func (s *Storage) RemoveRefresh(code string) error {
	fmt.Printf("RemoveRefresh: %s\n", code)
	delete(s.refresh, code)
	return nil
}

//getClients 取得人資系統OAuth Client資料表
func (s *Storage) getClients() []*oauthClients {
	mysql := MySQL{}.New()
	list := []*oauthClients{}
	err := mysql.Db.Find(&list).Error
	if err != nil {
		log.Fatal(err)
	}
	return list
}

//saveAuthorize2Redis 將驗証碼寫到Redis去
func (s *Storage) saveAuthorize2Redis(key, value string) bool {
	r := Redis{}.New()
	result := r.Set(key, value)
	if !result {
		return false
	}
	expire, err := strconv.Atoi(os.Getenv("EXPIRE_TIME"))
	if err != nil {
		log.Panicln(err)
		return false
	}
	r.SetExpire(key, expire)
	return true
}

//loadAuthorize2Redis 從Redis取得驗証碼
func (s *Storage) loadAuthorize2Redis(key string) string {
	r := Redis{}.New()
	result := r.Get(key)
	return result
}
