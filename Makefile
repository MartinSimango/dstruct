GO_PACKAGE:=github.com/MartinSimango/dstruct
TEST_FOLDER:=./tests
BENCHMARK_FOLDER:=./benchmarks

run:
	go run main/main.go

test: 
	go test ./...

test-versbose:
	go test -v ./...

test-cover:
	go test -v $(TEST_FOLDER)/dreflect_test -coverpkg $(GO_PACKAGE)/dreflect -coverprofile=$(TEST_FOLDER)/dreflect_coverage.out
	go test -v $(TEST_FOLDER)/dstruct_test -coverpkg $(GO_PACKAGE) -coverprofile=$(TEST_FOLDER)/dstruct_coverage.out

go tool cover -html=coverage.out


bench:
	go test -bench=.