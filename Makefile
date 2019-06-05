VERSION?=v0.0.2
PATH_BUILD=build/
FILE_COMMAND=lacuna
FILE_ARCH=linux_amd64

.PHONY: build clean version install run_show_debug

clean_build_and_test: version clean build run_show_debug

build:
	@$(GOPATH)/bin/goxc \
	  -bc="linux,amd64" \
	  -pv=$(VERSION) \
	  -d=$(PATH_BUILD) \
	  -build-ldflags "-X main.VERSION=$(VERSION)"
	sudo setcap cap_net_raw=+ep $(PATH_BUILD)$(VERSION)/$(FILE_ARCH)/$(FILE_COMMAND)

clean:
	@rm -rf ./build

version:
	@echo $(VERSION)

install:
	install -d -m 755 '$(HOME)/bin/'
	install $(PATH_BUILD)$(VERSION)/$(FILE_ARCH)/$(FILE_COMMAND) '$(HOME)/bin/$(FILE_COMMAND)'

run_show_debug:
	$(PATH_BUILD)$(VERSION)/$(FILE_ARCH)/$(FILE_COMMAND) show -d