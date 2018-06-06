package kftype

type Request struct {
	Option    uint8
	Done      chan error
	Namespace string
	Ingress   string
	Pod       string
}

type Response struct {
	Content string
}

const (
	IngressRoleAdd    = iota //0
	IngressRoleDelete        //1
)
