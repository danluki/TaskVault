
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
	rm -rf ui-dist
	rm -rf ui/build
	rm -rf ui/node_modules
	GOBIN=`pwd` go clean -i ./builtin/...
	GOBIN=`pwd` go clean

.PHONY: docs apidoc test ui updatetestcert
docs:
	# scripts/run doc --dir website/docs/cli
	
	# Build with docker while bun reach compatibility with docusaurs
	cd web; pnpm build --out-dir ../public
	ghp-import -p public

localtest:
	go test -v ./... | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''

updatetestcert:
	wget https://badssl.com/certs/badssl.com-client.p12 -q -O badssl.com-client.p12
	openssl pkcs12 -in badssl.com-client.p12 -nocerts -nodes -passin pass:badssl.com -out plugin/http/testdata/badssl.com-client-key-decrypted.pem
	openssl pkcs12 -in badssl.com-client.p12 -nokeys -passin pass:badssl.com -out plugin/http/testdata/badssl.com-client.pem
	rm badssl.com-client.p12

ui/node_modules: ui/package.json
	cd ui; bun install
	# touch the directory so Make understands it is up to date
	touch ui/node_modules

taskvault/ui-dist: ui/node_modules ui/public/* ui/src/* ui/src/*/*
	rm -rf taskvault/ui-dist
	cd ui; bun run build --out-dir ../taskvault/ui-dist

proto: types/taskvault.pb.go

types/%.pb.go: proto/%.proto
	protoc -I proto/ --go_out=types --go_opt=paths=source_relative --go-grpc_out=types --go-grpc_opt=paths=source_relative $<

client:
	oapi-codegen -package client website/static/openapi/openapi.yaml > client/client.go

ui: taskvault/ui-dist

main: taskvault/ui-dist types/taskvault.pb.go types/executor.pb.go *.go */*.go */*/*.go */*/*/*.go
	GOBIN=`pwd` go install ./builtin/...
	go mod tidy
	go build main.go