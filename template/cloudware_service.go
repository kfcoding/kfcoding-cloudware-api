package template

var CloudwareService = `
{
	"apiVersion": "v1",
	"kind": "Service",
	"metadata": {
		"name": "cloudware-3",
		"namespace": "kfcoding-alpha",
		"labels": {
			"app": "cloudware-3"
		}
	},
	"spec": {
		"ports": [
			{
				"port": 9800,
				"targetPort": 9800
			}
		],
		"selector": {
			"app": "cloudware-3"
		}
	}
}
`