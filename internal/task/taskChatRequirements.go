package task

import (
	"aurora/internal/chatgpt"
	"aurora/internal/tokens"
	"aurora/pkg/redis"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"
)

var errors []string

func CacheCreate() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	//取出所有的key
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println(time.Now(), tokens.CacheListKey)
			errors = []string{}
			preheatData, err := redis.Redis.LRange(tokens.CacheListKey, 1, 30)
			if len(preheatData) > 0 && err == nil {
				for _, v := range preheatData {
					go AddTaskRequireMent(v)
				}
			}
		}
		//break
	}
}

func AddTaskRequireMent(v string) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	if v != "" {
		ipSubsStr, err := url.QueryUnescape(v)
		if err != nil {
			log.Println("长度错误", err)
			return
		}
		ipSubs := strings.Split(ipSubsStr, "|")
		log.Println(ipSubs)

		if len(ipSubs) < 2 {
			log.Println("长度错误", v)
			return
		}

		key := "data:pow:" + ipSubs[0]
		if ok, _ := redis.Redis.Exists(key); !ok {
			log.Println("不存在")
			//config, err := send.GetGptPushConfig(v)

			//if err != nil {
			//	return
			//}

			if contains(errors, ipSubs[0]) {
				log.Println("已经不存在，不重新请求", v)
				return
			}

			//开始生成 token 有消息是5分钟
			//token, _, err := send.OpenaiSentinelChatRequirementsToken(config.Token, config.ProxyUrl, config.Id)
			//log.Println("token", token)
			if ipSubs[0] != "" && ipSubs[1] != "" {

				proxyUrl := tokens.Ipv6Set(ipSubs[1])
				turnStile, status, err := chatgpt.InitTurnStile(ipSubs[0], proxyUrl)
				if err != nil || status != 200 {
					return
				}

				if turnStile.TurnStileToken != "" {
					log.Println("正确", turnStile)
				}

				//log.Println("token, err", token, err)
				//if token.Token != "" && err == nil {
				//	//使用里面有两种
				//	var cacheToken v1.CacheToken
				//	if token.Token != "" {
				//		cacheToken.OpenaiSentinelChatRequirementsToken = token.Token
				//		cacheToken.OpenaiSentinelProofToken = token.OpenaiSentinelProofToken
				//	}
				//	jsonData, err := json.Marshal(cacheToken)
				//	if err != nil {
				//		return
				//	}
				//	a, b := redis.Redis.Set(key, jsonData, 300)
				//	log.Println(a, b, 300)
				//}
			} else {
				//本轮不在参与请求
				errors = append(errors, v)
			}
		}
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
