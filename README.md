# Syncra distributed core
The purpose of this repository is to educate how to build statefull services with Raft.
It's a nice way to get started with hashicorp Raft and Serf. To understand domain I can recommend
dkron repository(the base project for this simplified version) and Travis Jeffrey book "Distributed Services with Go".

So, it's just simple key value storage with eventual consistency with only one region support, in which you can put values like this
```sh
curl -X POST "http://localhost:8080/v1/storage" -H "Content-Type: application/json" -d '{"key": "test_key", "value": "test_value"}'
```
And get some with web ui or
```sh
curl http://localhost:8080/v1/storage                                                                                              
```

There is no Multi-Raft or multi regional support or distributed tx support and only few units and integrations test,
probably later this README will be updated with link to repsoitory to advanced version of this core. But for now I dunno how
to implement this to provide needed guarantees for this distributed system.

To test locally you can use any of compose files being here. For example
```sh
docker compose -f syncra-demo.yml up -d
```
Will up 3 nodes cluster with web ui on 8080,8081,8082. Default login and password is admin/admin.
