## How to run

```
$ docker build . -t wallex_test:prod
$ docker run -td wallex_test:prod
```

Now you can exec into the container and run the following commands

```
docker exec -it <container_id> ./vm_go
docker exec -it <container_id> ./vm_go_test -test.v
```
