test:
	@echo "Run project tests"
	@echo "------------------------------------------------------------"
	@go test -coverprofile=coverage.out -coverpkg=./... ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo ""
	@echo "Coverage output files Save!"
