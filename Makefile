export GO111MODULE=on

format: ## Format go code with goimports
	@go install github.com/rinchsan/gosimports/cmd/gosimports@latest
	@find . -name \*.go -exec gosimports -local github.com/coretrix/clockwork/ -w {} \;

#format-check: ## Check if the code is formatted
#	@go install -v golang.org/x/tools/cmd/goimports@latest
#	@for i in $$(goimports -l .); do echo "Code is not formatted run 'make format'" && exit 1; done

check: #format-check cyclo ## Linting and static analysis
	@if grep -r --include='*.go' -E "[^\/\/ ]+(fmt.Print|spew.Dump)"  *; then \
		echo "code contains fmt.Print* or spew.Dump function"; \
		exit 1; \
	fi

	@if test ! -e ./bin/golangci-lint; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.46.2; \
	fi
	@./bin/golangci-lint run --timeout 180s -E gosec -E stylecheck -E revive -E goimports -E whitespace

#static-check: format-check ## Linting and static analysis
#	@if grep -r --include='*.go' -E "fmt.Print|spew.Dump" *; then \
#		echo "code contains fmt.Print* or spew.Dump function"; \
#		exit 1; \
#	fi
#
#	@go install honnef.co/go/tools/cmd/staticcheck@latest;
#	@${GOPATH}/bin/staticcheck ./...;
#
#cyclo: ## Cyclomatic complexities analysis
#	@go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
#	@gocyclo -over 100 .

test: ## Run tests
	@go test -race -p 1 ./...
