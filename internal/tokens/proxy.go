package tokens

import (
	"aurora/pkg/redis"
	"aurora/pkg/setting"
	"bufio"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"net"
	"net/url"
	"os"
	"strings"
	"time"
)

var ips map[string]string

// var ids []string
var proxyUrl string
var proxyPrefix string
var Pingx string
var CacheListKey string

func Ping(c *gin.Context) {
	if Pingx == "1" {
		c.JSON(503, "server error")
	} else {
		c.JSON(200, "pong")
	}
}

func GetConfig(c *gin.Context) {
	c.JSON(200, gin.H{
		"PUID":         "",
		"PUIDTIME":     "",
		"UploadStatus": "500,401,422",
		"TokenUrl":     proxyUrl,
	})
}

func SetConfig(c *gin.Context) {
	var jsons struct {
		PUID         string `json:"puid"`
		TokenUrl     string `json:"token_url"`
		UploadStatus string `json:"upload_status"`
		VpsProxy     string `json:"vps_proxy"`
	}

	if err := c.BindJSON(&jsons); err != nil {
		return
	}

	if jsons.VpsProxy != "" {
		proxyUrl = jsons.VpsProxy
	}

	c.JSON(200, map[string]interface{}{
		"PUID":         "",
		"PUIDTIME":     "",
		"TokenUrl":     "",
		"UploadStatus": "",
		"VpsProxy":     proxyUrl,
	})
}

func init() {
	setting.Setup()
	redis.Setup()
	_ = godotenv.Load(".env")
	proxyUrl = os.Getenv("PROXY_URL")
	proxyPrefix = os.Getenv("PROXY_PREFIX")
	ipv64, err := readLines("ipv64.txt")
	if len(ipv64) == 0 || err != nil {
		log.Println("启动失败")
		return
	}
	CacheListKey = "data:tokens:list:" + os.Getenv("SERVER_PORT")

	ips = make(map[string]string)
	var ipsList []string
	for _, ipv64 := range ipv64 {
		UUID := uuid.NewString()
		//ids = append(ids, UUID)
		//ips.Store(UUID, ipv64)
		//log.Println(ids)
		if UUID != "" && ipv64 != "" {
			//ips[UUID] = ipv64
			ipsList = append(ipsList, url.QueryEscape(UUID+"|"+ipv64))
		}
	}

	if len(ipsList) > 0 {
		redis.Redis.LushS(CacheListKey, ipsList)
	}
}

func GetCacheList() (string, string, error) {
	//取出，存入
	ipSub, err := redis.Redis.LPop(CacheListKey)
	if err != nil || ipSub.(string) == "" {
		return "", "", errors.New("redis error")
	}

	ipSubStr, _ := url.QueryUnescape(ipSub.(string))
	redis.Redis.RPush(CacheListKey, ipSub.(string))
	p := strings.Split(ipSubStr, "|")

	if len(p) < 2 {
		return "", "", errors.New("redis error -1")
	}

	return p[0], p[1], nil
}

func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

//func ConfigProxy() (string, string) {
//	deviceId := randomString(ids)
//	if value, ok := ips[deviceId]; ok {
//		log.Println(value, ok)
//		return deviceId, value
//	}
//
//	return deviceId, ""
//}

// 随机
func randomString(strings []string) string {
	// 创建一个新的随机数生成器
	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)

	// 生成一个介于0和strings数组长度之间的随机数
	index := generator.Intn(len(strings))

	// 返回随机选中的字符串
	return strings[index]
}

// Ipv6Set 封装代理方式,每个/64 无限的 128地址 30 qps
func Ipv6Set(ipv6 string) string {
	if proxyUrl != "" {
		ipSub := strings.Split(proxyUrl, "|")
		if len(ipSub) == 2 {
			strIpv6, err := Ipv6New(ipSub[0] + ":" + ipv6)
			if err != nil {
				strIpv6 = ipSub[0] + ":" + ipv6
			}
			return "http://" + proxyPrefix + ":" + url.QueryEscape(strIpv6) + "@" + ipSub[1]
		}
	}

	return ""
}

// Ipv6New 随机生成
func Ipv6New(ipv6 string) (string, error) {
	// 定义子网前缀
	subnet := ipv6 + "/64"

	// 解析子网地址
	_, network, err := net.ParseCIDR(subnet)
	if err != nil {
		return "", err
	}

	// 获取子网的起始地址作为基地址
	baseIP := network.IP

	// 生成随机的后64位
	randomBytes := make([]byte, 8) // 64位 = 8字节
	_, err = rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// 将生成的随机数填充到基地址的后64位
	for i, b := range randomBytes {
		baseIP[8+i] = b
	}

	return baseIP.String(), nil
}
