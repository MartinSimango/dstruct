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

test-cover: test-dreflect test-dstruct

test-cover-verbose: test-dreflect-verbose test-dstruct-verbose

test-dreflect-verbose:
	go test -v $(DREFLECT_TEST)

test-dstruct-verbose:
	go test -v $(DSTRUCT_TEST)

test-dreflect:
	go test $(DREFLECT_TEST)

test-dstruct:
	go test $(DSTRUCT_TEST)

test-bench:
	go test $(TEST_FOLDER)/dstruct_test  -bench=.

# co
# go tool cover -html=coverage.out


bench:
	go test -bench=.


task-example:
	go run examples/task/main.go

### ACT ###
act:
	act -P ubuntu-latest=catthehacker/ubuntu:act-latest


### DAGGER ###
release-with-gpg-pass:
	@dagger call -m gopkg with-git-gpg-config \
		--gpg-key=env://GPG_KEY \
		--gpg-key-id=env://GPG_KEY_ID \
		--gpg-password=env://GPG_PASSPHRASE \
		--git-author-name "semantic-release-bot" \
		--git-author-email "shukomango@gmail.com" \
		release \
		--source=. \
		--dry-run=false \
		--token=env://GITHUB_TOKEN

release-with-gpg-no-pass:
	@dagger call -m gopkg with-git-gpg-config \
		--gpg-key=env://GPG_KEY \
		--gpg-key-id=env://GPG_KEY_ID \
		--git-author-name "semantic-release-bot" \
		--git-author-email "shukomango@gmail.com" \
		release \
		--source=. \
		--dry-run=false \
		--token=env://GITHUB_TOKEN


release:
	@dagger call -m gopkg release \
		--source=. \
		--dry-run=false \
		--token=env://GITHUB_TOKEN

release-dry-run:
	@dagger call -m gopkg release \
		--source=. \
		--dry-run=true \
		--token=env://GITHUB_TOKEN



