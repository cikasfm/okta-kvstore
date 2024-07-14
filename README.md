# Distributed Key-Value Store

## Summary

This project will try to implement and showcase development skills using GoLang

Author: [Zilvinas Vilutis](https://www.linkedin.com/in/zvilutis/)

## Running

To run the program you can use `make` commands:

### Testing

#### Unit Tests

```shell
make test
```

#### Manual Tests

Run the following commands once deployed

```shell
curl --location 'http://localhost:8080/key' \
--header 'Content-Type: application/json' \
--data '{
    "key1": "value1"
}'
```

### Running

Running a single instance

```shell
make run
```

OR

```shell
./main.sh
```