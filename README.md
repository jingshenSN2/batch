# Delay Batching

## Dummy worker

Start dummy worker.
```shell
go run worker/dummy_worker.go localhost:8089
```

Profile using python script.
```shell
python examples/dummy_client.py http://localhost:8089 1 2 3 5 10 50 100 300 1000
```


