package types

type RequestBody struct {
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
	Data  string
	Error string
}
