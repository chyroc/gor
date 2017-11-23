all: build

export GOBIN := $(CURDIR)/bin

build:
	$(BUILDENVVAR) go install -ldflags "-X main.BUILD_TIME=`date '+%Y-%m-%d_%I:%M:%S%p'` -X main.GIT_HASH=`git rev-parse HEAD`" github.com/Chyroc/gor

checkstyle:
	./.circleci/check_code_style.sh

test:
	gotest github.com/Chyroc/gor
	gotest github.com/Chyroc/gor/test

.PHONY: build checkstyle test
