.PHONY: all build push help

BIN_FILE=ad-exporter


linuxBuild:
	# default linux (centos,ubuntu,debian  other)
	@CGO_ENABLED=0 go build -o "${BIN_FILE}"

alpineBuild:
	# os alpine add: -tags natgo -installsuffix cgo
	@CGO_ENABLED=0 go build -a -installsuffix cgo -o "${BIN_FILE}"


build:  linuxBuild


push: 
	# .gitignore 如不生效: git rm -r --cached .
	@git add . ; git commit -m "`cat ./gitCommit`" ; git push

restart:
	@docker-compose down && docker-compose up -d



help:
	@echo "make build  编译生成二进制文件"
	@echo "make push git push 至代码仓库"


