## 云件控制器

1. 提供云件转发规则的增加及删除接口（直接修改traefik监听的etcd）

2. 云件保活，并在云件超时后调用restapi删除云件


POST    /api/cloudware/routing
DELETE  /api/cloudware/routing

POST    /api/traefik/routings
DELETE  /api/traefik/routings

POST    /api/cloudware/keepalive