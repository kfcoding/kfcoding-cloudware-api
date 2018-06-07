package configs

import (
	"os"
	"strconv"
	"log"
	"time"
)

var (
	ApiServerAddress = "http://api.kfcoding.com"
	ServerAddress    = "0.0.0.0:8080"
	PrefixAlive      = "/kfcoding/alpha/cloudware/ttl/"
	Namespace        = "kfcoding-alpha"
	IngressName      = "kfcoding-cloudware-ingress"
	WsAddrSuffix     = ".cloudware.kfcoding.com"
	QueueSize        = 1000
)

const (
	CloudWareTTL    = 60
	RequestTimeout  = 10 * time.Second
	EtcdDialTimeout = 5 * time.Second
)

func GetEtcdEndPoints() []string {
	return []string{"http://localhost:2379"}
}

func InitEnv() {
	var err error
	if ApiServerAddress = os.Getenv("ApiServerAddress"); "" == ApiServerAddress {
		ApiServerAddress = "http://api.kfcoding.com"
	}
	if ServerAddress = os.Getenv("ServerAddress"); "" == ServerAddress {
		ServerAddress = "0.0.0.0:8080"
	}
	if PrefixAlive = os.Getenv("PrefixAlive"); "" == PrefixAlive {
		PrefixAlive = "/kfcoding/alpha/cloudware/ttl/"
	}
	if Namespace = os.Getenv("Namespace"); "" == Namespace {
		Namespace = "kfcoding-alpha"
	}
	if IngressName = os.Getenv("IngressName"); "" == IngressName {
		IngressName = "kfcoding-cloudware-ingress"
	}
	if WsAddrSuffix = os.Getenv("WsAddrSuffix"); "" == WsAddrSuffix {
		WsAddrSuffix = ".cloudware.kfcoding.com"
	}

	if t := os.Getenv("QueueSize"); "" != t {
		if QueueSize, err = strconv.Atoi(t); nil != err {
			log.Fatal(err)
		}
	}
	log.Print(QueueSize)
	log.Print(ApiServerAddress)
	log.Print(ServerAddress)
	log.Print(PrefixAlive)
	log.Print(Namespace)
	log.Print(IngressName)
	log.Print(WsAddrSuffix)

}
