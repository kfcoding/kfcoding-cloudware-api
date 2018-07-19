package configs

import (
	"os"
	"log"
	"strings"
	"strconv"
	"time"
)

var (
	Namespace     = "kfcoding-alpha"
	ServerAddress = "0.0.0.0:8080"
	Token         = "kfcoding"
)
/**************** etcd config *************************************/
var (
	EtcdEndPoints  = []string{"http://localhost:2379"}
	EtcdUsername   = ""
	EtcdPassword   = ""
	RequestTimeout = 10 * time.Second
)
/**************** keep alive  *************************************/
/**************** KeeperPrefix/Version/TypeCloudware/name *********/
var KeeperTTL int64 = 60 // 保活时间
var (
	KeeperPrefix  = "/kfcoding/keepalive" // 保活前缀
	TypeCloudware = "1"                   // 云件类型
	Version       = "v1"                  // 版本
)
/**************** routing config **********************************/
var (
	PrefixTraefik = "/kfcoding/traefik/"
	WsAddrSuffix  = "cloudware.kfcoding.com"
)

func InitEnv() {

	if namespace := os.Getenv("Namespace"); namespace != "" {
		Namespace = namespace
	}
	if serverAddress := os.Getenv("ServerAddress"); serverAddress != "" {
		ServerAddress = serverAddress
	}
	if token := os.Getenv("Token"); token != "" {
		Token = token
	}

	/************************* etcd config ******************************/
	if etcdEndPoint := os.Getenv("EtcdEndPoints"); "" != etcdEndPoint {
		EtcdEndPoints = strings.Split(etcdEndPoint, ",")
	}
	if etcdUsername := os.Getenv("EtcdUsername"); "" != etcdUsername {
		EtcdUsername = etcdUsername
	}
	if etcdPassword := os.Getenv("EtcdPassword"); "" != etcdPassword {
		EtcdPassword = etcdPassword
	}

	/************************* keep alive config ************************/
	if ttl := os.Getenv("KeeperTTL"); "" != ttl {
		if t, err := strconv.ParseInt(ttl, 10, 64); nil != err {
			log.Fatal(err)
		} else {
			KeeperTTL = t
		}
	}
	if keeperPrefix := os.Getenv("KeeperPrefix"); keeperPrefix != "" {
		KeeperPrefix = keeperPrefix
	}
	if typeCloudware := os.Getenv("TypeCloudware"); typeCloudware != "" {
		TypeCloudware = typeCloudware
	}

	/*************************  routing config **************************/
	if prefixTraefik := os.Getenv("PrefixTraefik"); prefixTraefik != "" {
		PrefixTraefik = prefixTraefik
	}
	if wsAddrSuffix := os.Getenv("WsAddrSuffix"); wsAddrSuffix != "" {
		WsAddrSuffix = wsAddrSuffix
	}

	log.Print("ServerAddress:  ", ServerAddress)
	log.Print("Namespace:      ", Namespace)
	log.Print("Token:          ", Token)
	log.Print("KeeperTTL:      ", KeeperTTL)
	log.Print("KeeperPrefix:   ", KeeperPrefix)
	log.Print("TypeCloudware:  ", TypeCloudware)
	log.Print("Version:        ", Version)
	log.Print("EtcdEndPoints:  ", EtcdEndPoints)
	log.Print("PrefixTraefik:  ", PrefixTraefik)
	log.Print("WsAddrSuffix:   ", WsAddrSuffix)

	/*
	if ServerAddress =; "" == ServerAddress {
		ServerAddress = "0.0.0.0:8080"
	}
	if Namespace = os.Getenv("Namespace"); "" == Namespace {
		Namespace = "kfcoding-alpha"
	}
	if Token = os.Getenv("Token"); "" == Token {
		Token = "Bearer ad3efe453a786f036a946015feff19f78a80192f462ea1d56e3d89e8c4f5d833"
	}

	if PrefixTraefik = os.Getenv("PrefixTraefik"); "" == PrefixTraefik {
		PrefixTraefik = "/kfcoding/cloudware/traefik/"
	}
	if WsAddrSuffix = os.Getenv("WsAddrSuffix"); "" == WsAddrSuffix {
		WsAddrSuffix = "cloudware.kfcoding.com"
	}

	if EtcdEndPoint := os.Getenv("EtcdEndPoints"); "" == EtcdEndPoint {
		EtcdEndPoints = []string{"http://localhost:2379"}
		//EtcdEndPoints = []string{"http://10.99.139.170:2379"}
	} else {
		EtcdEndPoints = strings.Split(EtcdEndPoint, ",")
	}
	if ttl := os.Getenv("CloudWareTTL"); "" != ttl {
		if t, err := strconv.ParseInt(ttl, 10, 64); nil != err {
			log.Fatal(err)
		} else {
			CloudWareTTL = t
		}
	} else {
		CloudWareTTL = 60
	}
	*/

}
