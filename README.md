# Syncra distributed core
The purpose of this repository is to educate how to build statefull services with Raft.
It's a nice way to get started with hashicorp Raft and Serf. To understand domain I can recommend
dkron repository(the base project for this simplified version) and Travis Jeffrey book "Distributed Services with Go".

## Demo
Demo is deployed on syncra[1-3].danluki.ru to get access to test web ui. Default login and password is admin/admin.

To test locally you can use any of compose files being here. For example
```
docker compose -f syncra-demo.yml up -d
```
Will up 3 nodes cluster with web ui on 8080,8081,8082
