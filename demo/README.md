This is simple demo deployment with high availability and node scaling. Only docker required

```
    docker compose -f docker-compose.yml up
```

After looking for http://localhost:8080/ui do

```
    docker compose -f compose2.yml up
```

You will see Syncra deployment with 5 nodes.