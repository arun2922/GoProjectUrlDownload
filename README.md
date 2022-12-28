# GoProjectUrlDownload

## Install
```bash
git clone https://github.com/arun2922/GoProjectUrlDownload.git
go get .
```

## Run

In one of the terminals run (this acts as server which will be running on localhost:7771)
`go run main.go`

In another terminal enter below command for POST request

`curl -H "Content-Type: application/json" -d "{\"Uri\":\"https://google.com\",\"RetryLimit\":3}" -X POST http://localhost:7771/pagesource`

