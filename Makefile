GO_PACKAGE:=github.com/MartinSimango/dstruct
TEST_FOLDER:=./tests
BENCHMARK_FOLDER:=./benchmarks


TESTS:=$(TEST_FOLDER)/dreflect_test $(TEST_FOLDER)/dstruct_test

DREFLECT_TEST:=$(TEST_FOLDER)/dreflect_test -coverpkg $(GO_PACKAGE)/dreflect -coverprofile=$(TEST_FOLDER)/dreflect_test/coverage.txt
DSTRUCT_TEST:=$(TEST_FOLDER)/dstruct_test -coverpkg $(GO_PACKAGE) -coverprofile=$(TEST_FOLDER)/dstruct_test/coverage.txt

run:
	go run main/main.go

test: 
	go test ./...

test-verbose:
	go test -v ./...

test-cover:
	go test $(DREFLECT_TEST)
	go test $(DSTRUCT_TEST)

test-cover-verbose:
	go test -v $(DREFLECT_TEST)
	go test -v $(DSTRUCT_TEST)

# co
# go tool cover -html=coverage.out


bench:
	go test -bench=.