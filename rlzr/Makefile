ifeq ($(uname),Linux)
	MAKE_LZR_MSG := "Linux OS detected. Proceed."
	MAKE_CLEAN_MSG := "Don't forget to delete iptables rule using:\nsudo iptables -L --line-numbers && sudo iptables -D OUTPUT \TK""
else
	MAKE_LZR_MSG := "WARNING: Non-Linux OS detected. You cannot run rlzr executable, but can run tests and build internally."
	MAKE_CLEAN_MSG := ""
endif

GO_FILES = $(shell find . -type f -name '*.go')

all: lzr
	sudo iptables -A OUTPUT -p tcp --tcp-flags RST RST -s $(source-ip) -j DROP

lzr: $(GO_FILES)
	@echo $(MAKE_LZR_MSG)
	cd cmd/rlzr && go build -buildvcs=false && cd ../..
	rm -f rlzr
	ln -s cmd/rlzr/rlzr rlzr

clean:
	cd cmd/rlzr && go clean && rm -f rlzr && cd ../..
	rm -f rlzr
	@echo $(MAKE_CLEAN_MSG)
