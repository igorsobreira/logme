
ifeq ($(package),)
	PACKAGE=./...
else
	PACKAGE=github.com/igorsobreira/logme/$(package)
endif

# Provide custom test args with:
#
#  make test args='-v'
#
# test a single package with:
#
#  make test package=handlers
#
test:
	go test -race -i ./...
	go test -race $(args) $(PACKAGE)

fmt:
	@go fmt ./...
