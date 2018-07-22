package types

//type RequestBody struct {
//	Name string
//	URL  string
//	Rule string
//}

type RoutingBody struct {
	Name string
	URL  string
	Rule string
}

type KeeperBody struct {
	Name string
}

type CloudwareBody struct {
	Image string
}

type ResponseBody struct {
	Name  string
	Data  string
	Error string
}
