package kftype

type Request struct {
	Option    uint8
	Done      chan error
	Namespace string
	Ingress   string
	Pod       string
}

const (
	INGRESS_RULE_ADD    = iota //0
	INGRESS_RULE_DELETE        //1
)
