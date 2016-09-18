alerts1='[
  {
        "labels": {
            "__name__": "node_memory_Active",
            "alertname": "MemoryActiveOutB",
            "instance": "10.10.12.18:9100",
            "job": "node",
            "monitor": "codelab-monitor",
            "team":"cloud",
            "asdasfasfasfasfasfasfasfasfasfasfafs":"asfasfasfasfasfasfasfasfasfasfasfasfafs"
        },
        "annotations": {
            "description": "10.10.12.18:9100 of node memory active out "
        },
        "startsAt": "2016-08-29T07:39:52.806Z",
        "endsAt": "0001-01-01T00:00:00Z",
        "generatorURL": "http://987a66675dda:9090/graph#%5B%7B%22expr%22%3A%22node_memory_Active%7Binstance%3D%5C%2210.10.12.18%3A9100%5C%22%7D%20%3E%202000000%22%2C%22tab%22%3A0%7D%5D"
    },
    {
        "labels": {
            "__name__": "node_memory_Active",
            "alertname": "MemoryActiveOutB",
            "instance": "10.10.12.18:9100",
            "job": "node",
            "monitor": "codelab-monitor",
            "team":"cloud"
        },
        "annotations": {
            "description": "10.10.12.18:9100 of node memory active out "
        },
        "startsAt": "2016-08-29T07:39:52.806Z",
        "endsAt": "0001-01-01T00:00:00Z",
        "generatorURL": "http://987a66675dda:9090/graph#%5B%7B%22expr%22%3A%22node_memory_Active%7Binstance%3D%5C%2210.10.12.18%3A9100%5C%22%7D%20%3E%202000000%22%2C%22tab%22%3A0%7D%5D"
    },
    {
        "labels": {
            "__name__": "node_memory_Active",
            "alertname": "MemoryActiveOutB",
            "instance": "10.10.12.18:9100",
            "job": "node",
            "monitor": "codelab-monitor",
            "team":"cloud"
        },
        "annotations": {
            "description": "10.10.12.18:9100 of node memory active out "
        },
        "startsAt": "2016-08-29T07:39:52.806Z",
        "endsAt": "0001-01-01T00:00:00Z",
        "generatorURL": "http://987a66675dda:9090/graph#%5B%7B%22expr%22%3A%22node_memory_Active%7Binstance%3D%5C%2210.10.12.18%3A9100%5C%22%7D%20%3E%202000000%22%2C%22tab%22%3A0%7D%5D"
    }

]'
curl -XPOST -d"$alerts1"  http://localhost:8888/api/v1/alerts
curl -XPOST -d"$alerts1" http://localhost:8888/api/v1/alerts
curl -XPOST -d"$alerts1" http://localhost:8888/api/v1/alerts