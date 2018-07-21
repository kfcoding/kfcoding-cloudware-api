package configs

import (
	"os"
	"log"
	"strings"
	"strconv"
	"time"
)

var (
	ServerAddress = "0.0.0.0:8080"
	Token         = ""
	Version       = "v1"
	Namespace     = "kfcoding-alpha"
)

// etcd config
var (
	EtcdEndPoints  = []string{"http://localhost:2379"}
	EtcdUsername   = ""
	EtcdPassword   = ""
	RequestTimeout = 10 * time.Second
)

// keep alive  KeeperPrefix/Version/name
var (
	KeeperPrefix = "/kfcoding/keepalive/cloudware"
	KeeperTTL    = 60
)

// routing config
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
	if version := os.Getenv("Version"); version != "" {
		Version = version
	}

	// etcd config
	if etcdEndPoint := os.Getenv("EtcdEndPoints"); "" != etcdEndPoint {
		EtcdEndPoints = strings.Split(etcdEndPoint, ",")
	}
	if etcdUsername := os.Getenv("EtcdUsername"); "" != etcdUsername {
		EtcdUsername = etcdUsername
	}
	if etcdPassword := os.Getenv("EtcdPassword"); "" != etcdPassword {
		EtcdPassword = etcdPassword
	}

	// keep alive config
	if ttl := os.Getenv("KeeperTTL"); "" != ttl {
		if t, err := strconv.ParseInt(ttl, 10, 64); nil != err {
			log.Fatal(err)
		} else {
			KeeperTTL = int(t)
		}
	}
	if keeperPrefix := os.Getenv("KeeperPrefix"); keeperPrefix != "" {
		KeeperPrefix = keeperPrefix
	}

	// routing config
	if prefixTraefik := os.Getenv("PrefixTraefik"); prefixTraefik != "" {
		PrefixTraefik = prefixTraefik
	}
	if wsAddrSuffix := os.Getenv("WsAddrSuffix"); wsAddrSuffix != "" {
		WsAddrSuffix = wsAddrSuffix
	}

	log.Print("ServerAddress:  ", ServerAddress)
	log.Print("Token:          ", Token)
	log.Print("Version:        ", Version)
	log.Print("Version:        ", Version)

	log.Print("EtcdEndPoints:  ", EtcdEndPoints)
	log.Print("EtcdUsername:   ", EtcdUsername)
	log.Print("EtcdPassword:   ", EtcdPassword)

	log.Print("KeeperTTL:      ", KeeperTTL)
	log.Print("KeeperPrefix:   ", KeeperPrefix)

	log.Print("PrefixTraefik:  ", PrefixTraefik)
	log.Print("WsAddrSuffix:   ", WsAddrSuffix)

}
