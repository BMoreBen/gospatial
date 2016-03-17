import json
import requests

req = requests.post("http://localhost:8888/api/v1/layer")
res = json.loads(req.json())
ds = res["datasource"]
print(res)

req = requests.get("http://localhost:8888/api/v1/layer/" + ds)
res = req.json()
print(res)

payload = {
	"geometry": {
		"type": "Point",
		"coordinates": [10,-10]
	},
	"properties": {
		"name": "test point 1"
	}
}

print("POST FEATURE")
req = requests.post("http://localhost:8888/api/v1/layer/" + ds + "/feature", data=json.dumps(payload))
print(req.json())

print()
print("GET FEATURE")
req = requests.get("http://localhost:8888/api/v1/layer/" + ds + "/feature/0")
print(req.json())

print()
print("DELETE FEATURE")
req = requests.delete("http://localhost:8888/api/v1/layer/" + ds)
print(req.json())
