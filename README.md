# tesla_automate_charge
Automatically control Tesla vehicle charging rate based on Solar Panel Generation

Currently, only support fetching Solar Data from Sense Monitoring

## Build
Build directly from the repo
```text
go build -trimpath -ldflags "-w -s"             . 
```

Build from docker
```text
# Linux
docker run --rm -v $(PWD):/go/src/builder -w /go/src/builder -e GOOS=linux -e GOARCH=amd64 golang:1.17 go build -trimpath -ldflags "-w -s"

# Mac Intel
docker run --rm -v $(PWD):/go/src/builder -w /go/src/builder -e GOOS=darwin -e GOARCH=amd64 golang:1.17 go build -trimpath -ldflags "-w -s"  

# Mac M1
docker run --rm -v $(PWD):/go/src/builder -w /go/src/builder -e GOOS=darwin -e GOARCH=arm64 golang:1.17 go build -trimpath -ldflags "-w -s"

# Windows docker run -v $(PWD):/go/src/builder -w /go/src/builder -e GOOS=windows -e GOARCH=amd64 golang:1.17 go build -trimpath -ldflags "-w -s"
docker run --rm -v $(PWD):/go/src/builder -w /go/src/builder -e GOOS=windows -e GOARCH=amd64 golang:1.17 go build -trimpath -ldflags "-w -s"
```

Build all platforms
```text
docker run --rm -v $(PWD):/go/src/builder -w /go/src/builder -e GOOS=windows -e GOARCH=amd64 golang:1.17 bash build.sh
```


## Usage
copy configs.yml.template to configs.yml

update configurations in configs.yml

Start the program
```bash
./tesla_automated_charge_control
```
