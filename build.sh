apt-get update -y
apt-get install zip -y

GOOS=windows
GOARCH=amd64
F_NAME=tesla_automated_charge_control-${GOOS}-${GOARCH}.exe
go build -trimpath -ldflags "-w -s" -o $F_NAME
zip $F_NAME.zip $F_NAME configs.yml.template
rm -rf $F_NAME

GOOS=darwin
GOARCH=amd64
F_NAME=tesla_automated_charge_control-${GOOS}-${GOARCH}
go build -trimpath -ldflags "-w -s" -o $F_NAME
zip $F_NAME.zip $F_NAME configs.yml.template
rm -rf $F_NAME

GOOS=darwin
GOARCH=arm64
F_NAME=tesla_automated_charge_control-${GOOS}-${GOARCH}
go build -trimpath -ldflags "-w -s" -o $F_NAME
zip $F_NAME.zip $F_NAME configs.yml.template
rm -rf $F_NAME

GOOS=linux
GOARCH=amd64
F_NAME=tesla_automated_charge_control-${GOOS}-${GOARCH}
go build -trimpath -ldflags "-w -s" -o $F_NAME
zip $F_NAME.zip $F_NAME configs.yml.template
rm -rf $F_NAME
