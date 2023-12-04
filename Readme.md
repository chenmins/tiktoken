## 测试说明

```
docker run --rm -it -p 8888:8888 chenmins/tiktoken

docker run -d -name tiktoken  -p 8888:8888 chenmins/tiktoken

```
 
## 测试用例

```
curl --location --request POST 'http://tx.chenmin.org:8888' \
--header 'Content-Type: application/json' \
--data-raw '{
    "latencies": {
        "proxy": 130,
        "request": 342,
        "kong": 212
    },
    "client_ip": "219.143.240.71",
    "request": {
        "size": 820,
        "body": "{\n    \"model\": \"gpt-3.5-turbo\",\n    \"stream\":true,\n    \"messages\": [\n      {\n        \"role\": \"system\",\n        \"content\": \"You are a helpful assistant.\"\n      },\n      {\n        \"role\": \"user\",\n       \"content\": \"1+1=2 对吗，如果正确，请仅仅回答yes?\"\n      } ,\n      {\n        \"role\": \"assistant\",\n       \"content\": \"yes\"\n      } ,\n      {\n        \"role\": \"user\",\n       \"content\": \"1+2=3 对吗，如果正确，请仅仅回答yes?\"\n      } \n    ]\n  }",
        "headers": {
            "user-agent": "PostmanRuntime/7.26.10",
            "authorization": "REDACTED",
            "content-type": "application/json",
            "content-length": "465",
            "connection": "keep-alive",
            "host": "tx.chenmin.org:8000",
            "accept": "*/*",
            "accept-encoding": "gzip, deflate, br",
            "postman-token": "ec21d3df-5af5-4663-9a1c-9d1043e6ff2f"
        },
        "method": "POST",
        "uri": "/v1/chat/completions",
        "url": "http://tx.chenmin.org:8000/v1/chat/completions",
        "querystring": {},
        "id": "7fd952507c095c09519a9220da452140"
    },
    "upstream_uri": "/v1/chat/completions",
    "started_at": 1701175944461,
    "service": {
        "enabled": true,
        "connect_timeout": 60000,
        "read_timeout": 60000,
        "host": "esb.gt.cn",
        "name": "esb",
        "id": "c79d01b7-86a0-498d-b461-7157c35685ed",
        "created_at": 1700223105,
        "updated_at": 1700223105,
        "ws_id": "f676f6e0-e2f6-4c59-8331-c1840a91c6a5",
        "protocol": "https",
        "port": 443,
        "retries": 5,
        "write_timeout": 60000
    },
    "route": {
        "request_buffering": true,
        "response_buffering": true,
        "regex_priority": 0,
        "https_redirect_status_code": 426,
        "name": "openai",
        "id": "92d2ee71-5979-4c85-af5b-26f7d771c380",
        "updated_at": 1700225108,
        "ws_id": "f676f6e0-e2f6-4c59-8331-c1840a91c6a5",
        "path_handling": "v0",
        "strip_path": false,
        "preserve_host": false,
        "tags": [],
        "service": {
            "id": "c79d01b7-86a0-498d-b461-7157c35685ed"
        },
        "protocols": [
            "http",
            "https"
        ],
        "paths": [
            "/v1"
        ],
        "created_at": 1700225021
    },
    "upstream_status": "200",
    "response": {
        "body": "data: {\"model\":\"gpt-3.5-turbo\",\"object\":\"chat.completion.chunk\",\"choices\":[{\"index\":0,\"delta\":{\"role\":\"assistant\"},\"finish_reason\":null}]}\r\n\r\ndata: {\"model\":\"gpt-3.5-turbo\",\"object\":\"chat.completion.chunk\",\"choices\":[{\"index\":0,\"delta\":{\"content\":\"\\n\"},\"finish_reason\":null}]}\r\n\r\ndata: {\"model\":\"gpt-3.5-turbo\",\"object\":\"chat.completion.chunk\",\"choices\":[{\"index\":0,\"delta\":{\"content\":\" yes\"},\"finish_reason\":null}]}\r\n\r\ndata: {\"model\":\"gpt-3.5-turbo\",\"object\":\"chat.completion.chunk\",\"choices\":[{\"index\":0,\"delta\":{\"content\":null},\"finish_reason\":null}]}\r\n\r\ndata: {\"model\":\"gpt-3.5-turbo\",\"object\":\"chat.completion.chunk\",\"choices\":[{\"index\":0,\"delta\":{},\"finish_reason\":\"stop\"}]}\r\n\r\ndata: [DONE]\r\n\r\n",
        "status": 200,
        "headers": {
            "ratelimit-reset": "36",
            "via": "kong/3.5.0",
            "x-kong-request-id": "7fd952507c095c09519a9220da452140",
            "content-type": "text/event-stream; charset=utf-8",
            "transfer-encoding": "chunked",
            "date": "Tue, 28 Nov 2023 12:52:20 GMT",
            "ratelimit-remaining": "4",
            "x-kong-proxy-latency": "0",
            "cache-control": "no-cache",
            "x-kong-upstream-latency": "130",
            "server": "uvicorn",
            "x-ratelimit-limit-minute": "5",
            "ratelimit-limit": "5",
            "x-ratelimit-remaining-minute": "4",
            "connection": "close"
        },
        "size": 1162
    },
    "tries": [
        {
            "balancer_start_ns": 1.7011759444617e+18,
            "ip": "58.144.220.44",
            "port": 443,
            "balancer_latency_ns": 16640,
            "balancer_latency": 0,
            "balancer_start": 1701175944461
        }
    ]
}'

```