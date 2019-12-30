#!/bin/bash
set -x

curl -si 127.0.0.1:8080/api/fetcher -X POST -d '{"url":"https://httpbin.org/range/5","interval":5}'
curl -si 127.0.0.1:8080/api/fetcher -X POST -d '{"url":"https://httpbin.org/range/50","interval":50}'i
curl -si 127.0.0.1:8080/api/fetcher -X POST -d '{"url":"https://httpbin.org/range/120","interval":60}'
curl -si 127.0.0.1:8080/api/fetcher -X POST -d '{"url":"https://httpbin.org/range/1000","interval":60}'

curl -s 127.0.0.1:8080/api/fetcher | python -m json.tool

