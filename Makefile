# Developer shortcuts. The filter itself needs only "go build".

.PHONY: build test docs clean

# build compiles doxygo into bin/.
build:
	go build -o bin/doxygo .

# test runs every test.
test:
	go vet ./...
	go test ./...

# docs builds doxygo's own documentation into build/docs, with doxygo.
docs: build
	mkdir -p build/docs
	DOXYGO=$(CURDIR)/bin/doxygo doxygen Doxyfile

# clean removes build output.
clean:
	rm -rf bin build
