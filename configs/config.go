package configs

import (
	"os"
	"strconv"
	"log"
	"time"
	"strings"
)

var (
	ApiServerAddress = "http://api.kfcoding.com"
	ServerAddress    = "0.0.0.0:8080"
	Namespace        = "kfcoding-alpha"
	Token            = "Bearer ad3efe453a786f036a946015feff19f78a80192f462ea1d56e3d89e8c4f5d833"

	PrefixAlive   = "/kfcoding/cloudware/ttl/"
	PrefixTraefik = "/kfcoding/cloudware/traefik/"
	WsAddrSuffix  = "cloudware.kfcoding.com"
	EtcdEndPoints = []string{}
	CloudWareTTL  int64
)

const (
	RequestTimeout = 10 * time.Second
)

func InitEnv() {
	if ApiServerAddress = os.Getenv("ApiServerAddress"); "" == ApiServerAddress {
		ApiServerAddress = "http://api.kfcoding.com"
	}
	if ServerAddress = os.Getenv("ServerAddress"); "" == ServerAddress {
		ServerAddress = "0.0.0.0:8080"
	}
	if Namespace = os.Getenv("Namespace"); "" == Namespace {
		Namespace = "kfcoding-alpha"
	}
	if Token = os.Getenv("Token"); "" == Token {
		Token = "Bearer ad3efe453a786f036a946015feff19f78a80192f462ea1d56e3d89e8c4f5d833"
	}

	if PrefixAlive = os.Getenv("PrefixAlive"); "" == PrefixAlive {
		PrefixAlive = "/kfcoding/cloudware/ttl/"
	}
	if PrefixTraefik = os.Getenv("PrefixTraefik"); "" == PrefixTraefik {
		PrefixTraefik = "/kfcoding/cloudware/traefik/"
	}
	if WsAddrSuffix = os.Getenv("WsAddrSuffix"); "" == WsAddrSuffix {
		WsAddrSuffix = "cloudware.kfcoding.com"
	}

	if EtcdEndPoint := os.Getenv("EtcdEndPoints"); "" == EtcdEndPoint {
		EtcdEndPoints = []string{"http://etcd." + Namespace + ".svc.cluster.local:2379"}
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

	log.Print("ApiServerAddress: ", ApiServerAddress)
	log.Print("ServerAddress: ", ServerAddress)
	log.Print("Namespace: ", Namespace)
	log.Print("Token: ", Token)
	log.Print("PrefixAlive: ", PrefixAlive)
	log.Print("PrefixTraefik: ", PrefixTraefik)
	log.Print("WsAddrSuffix: ", WsAddrSuffix)
	log.Print("EtcdEndPoints: ", EtcdEndPoints)
	log.Print("CloudWareTTL: ", CloudWareTTL)

}
