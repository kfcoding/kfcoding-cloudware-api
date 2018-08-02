## 云件控制器(暂不使用)

- 根据镜像名创建云件
- 云件保活
- 添加路由
- 删除路由

1. build

```
go build -o controller main.go

build using docker golang

docker run -it -v /Users/wsl/go/src:/go/src golang:1.10.3-alpine3.8 sh

cd src/github.com/kfcoding-cloudware-controller/ && go build -o controller main.go && exit

scp controller root@worker:/home/kfcoding-cloudware-controller

cd /home/kfcoding-cloudware-controller && \
docker build -t daocloud.io/shaoling/kfcoding-cloudware-controller:v2.1 .

```

2. 创建云件

```
POST    /cloudware/

Header
    Content-Type: application/json
    Token ""
Body
    {
        "Image": string
    }
Response
    {
        Data:   string
        Name:   string
        Error:  string
    }
```

3. 云件保活

```
POST    /keep/cloudware

Header
    Content-Type: application/json

Body
    {
        "Name":"cloudware-example"
    }
Response
    {
        Data:   string
        Name:   string
        Error:  string
    }
```

4. 添加路由
```
POST    /routing/add
Header
    Content-Type: application/json
    Token: ""
Body
    {
        "Name": string
        "URL":  string
        "Rule": string
    }
ResponseBody
    {
        Data:   string
        Name:   string
        Error:  string
    }
```

5. 添加多条路由
```
POST    /routing/adds
Header
    Content-Type: application/json
    Token: ""
Body
    [
        {
            "Name": string
            "URL":  string
            "Rule": string
        },
        {
            "Name": string
            "URL":  string
            "Rule": string
        }
    ]
ResponseBody
    {
        Data:   string
        Name:   string
        Error:  string
    }
```

6. 删除路由
```
POST    /routing/delete
Header
    Content-Type: application/json
    Token: ""
Body
    {
        "Name": string
    }
ResponseBody
    {
        Data:   string
        Name:   string
        Error:  string
    }
```

7. 删除多条路由
```
POST    /routing/deletes
Header
    Content-Type: application/json
    Token: ""
Body
    [
           {
               "Name": string
           },
           {
               "Name": string
           }
    ]
ResponseBody
    {
        Data:   string
        Name:   string
        Error:  string
    }
```
