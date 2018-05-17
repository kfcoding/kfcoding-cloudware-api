package kftype

type Request struct {
	Option    uint8
	Done  chan string
	Namespace string
	Ingress   string
	Pod       string
}

const (
	INGRESS_RULE_ADD    = iota
	INGRESS_RULE_DELETE
)
