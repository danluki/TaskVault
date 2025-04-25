
.PHONY: clean
clean:
	rm -f main
	rm -f *_SHA256SUMS
	rm -f taskvault-*
	rm -rf build/*
	rm -rf builder/skel/*
	rm -f *.deb
	rm -f *.rpm
	rm -f *.tar.gz
	rm -rf tmp
	rm -rf node-ui/dist
	rm -rf node-ui/node_modules
	GOBIN=`pwd` go clean -i ./builtin/...
	GOBIN=`pwd` go clean


node-ui/node_modules: node-ui/package.json
	cd node-ui; bun install
	touch node-ui/node_modules

taskvault/ui-dist: node-ui/node_modules node-ui/public/* node-ui/src/* node-ui/src/*/*
	rm -rf taskvault/ui-dist
	cd node-ui; bun run build --out-dir ../taskvault/ui-dist

proto: pkg/types

types/%.pb.go: proto/%.proto
	protoc -I proto/ --go_out=types --go_opt=paths=source_relative --go-grpc_out=types --go-grpc_opt=paths=source_relative $<

ui: taskvault/ui-dist

main: taskvault/ui-dist pkg/types  *.go */*.go */*/*.go
	go mod tidy
	go build main.go

dev:
	go run main.go agent --bootstrap=true