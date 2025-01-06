BUILD_NAME:=logViewerAgent

SOURCE:=*.go

.DEFAULT_GOAL := default

.PHONY:macarm64 win_amd64 linux_amd64 clean pre_build post_build

build_mac_arm64:clean pre_build mac_arm64 post_build

build_win_amd64:clean pre_build win_amd64 post_build

build_linux_amd64:clean pre_build linux_amd64 post_build

mac_arm64: deps
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-w -s" -o ${BUILD_NAME} ${SOURCE}

win_amd64: deps
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o ${BUILD_NAME}.exe  ${SOURCE}

linux_amd64:deps
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ${BUILD_NAME} ${SOURCE}

pre_build:
	@echo "===start build==="
	@mkdir -pv dist/${BUILD_NAME}/log

post_build:
	@cp -r {conf,storage,README.MD,startup.sh} dist/${BUILD_NAME}
	@mv ${BUILD_NAME}* dist/${BUILD_NAME}
	@cd dist && tar -zcf ${BUILD_NAME}.tar.gz ${BUILD_NAME}
	@rm -rf dist/${BUILD_NAME}
	@echo "===end build==="
	@echo
	@echo "+++info+++"
	@echo " target file is \"dist/${BUILD_NAME}.tar.gz\""
	@echo
	@echo " \"README.MD\" shows how to deploy the service"
	@echo

deps:
	@go vet
	go mod tidy

clean:
	@rm -rf dist
	@go clean

default:
	@echo "  options:"
	@echo "	make build_mac_arm64 : build mac arm64 target file"
	@echo "	make build_win_amd64 : build windows amd64 target file"
	@echo "	make build_linux_amd64 : build linux amd64 target file"
	@echo
