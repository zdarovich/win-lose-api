# win-lose-api

# Run
```shell script
git clone https://github.com/zdarovich/win-lose-api
cd win-lose-api
docker-compose up
```

# Run test client
Test client calls /your_url endpoint with random data for 1 minute with 5 seconds interval.
```shell script
go run cmd/client/main.go
```

# Create valid transaction
```shell script
curl --location --request POST 'http://127.0.0.1:8081/your_url' \
--header 'Content-Type: application/json' \
--data-raw '
            {
               "state": "lost",
               "amount": "10.5",
               "transactionId": "beb4ad51-702d-4372-a7af-07801f430df2"
            }
'
```
