package libs

import (
	"log"
	"os"

	"github.com/garyburd/redigo/redis"
)

//Redis 相關設定
type Redis struct {
	client redis.Conn
}

//New 建構式
func (r Redis) New() *Redis {
	protocol := os.Getenv("REDIS_PROTOCOL")
	host := os.Getenv("REDIS_HOST")
	var err error
	r.client, err = redis.Dial(protocol, host)
	if err != nil {
		log.Panicln(err)
	}
	return &r
}

//SetExpire 設定資料過期時間
func (r *Redis) SetExpire(key string, expire int) {
	_, err := r.client.Do("EXPIRE", key)
	if err != nil {
		log.Panicln(err)
	}
	defer r.client.Close()
}

//Get 取得Redis資料
func (r *Redis) Get(key string) string {
	value, err := redis.String(r.client.Do("get", key))
	if err != nil {
		log.Panicln(err)
	}
	defer r.client.Close()
	return value
}

//Set 設定Redis資料
func (r *Redis) Set(key, value string) bool {
	_, err := r.client.Do("set", key, value)
	if err != nil {
		log.Panicln(err)
		return false
	}
	defer r.client.Close()
	return true
}
