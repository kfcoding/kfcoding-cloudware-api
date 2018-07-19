package template

var (
	CloudwarePod = `
{
	"apiVersion": "v1",
	"kind": "Pod",
	"metadata": {
		"name": "cloudware-3",
		"namespace": "kfcoding-alpha",
		"labels": {
			"app": "cloudware-3"
		}
	},
	"spec": {
		"containers": [
			{
				"name": "xorg",
				"command": [
					"Xorg"
				],
				"image": "daocloud.io/shaoling/kfcoding-xorg:master-094594c",
				"volumeMounts": [
					{
						"name": "app-tmp",
						"mountPath": "/tmp"
					}
				]
			},
			{
				"name": "pulsar",
				"image": "daocloud.io/shaoling/kfcoding-rstudio-latest:master",
				"ports": [
					{
						"containerPort": 9800
					}
				],
				"volumeMounts": [
					{
						"name": "app-tmp",
						"mountPath": "/tmp"
					}
				]
			}
		],
		"volumes": [
			{
				"name": "app-tmp",
				"emptyDir": {}
			}
		]
	}
}
`
)
