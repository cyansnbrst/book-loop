.PHONY: cover
cover:
	go test -short -count=1 -race -coverprofile=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

.PHONY: gen
gen:
	mockgen -source=internal/books/pg_repository.go -destination=internal/books/mock/mock_pg_repository.go
	mockgen -source=internal/books/usecase.go -destination=internal/books/mock/mock_usecase.go
 