
GO_CGO_CFLAGS := CGO_CFLAGS='-O1 -I$(CURDIR)/capstone/include'
GO_CGO_LDFLAGS := CGO_LDFLAGS='-O1 -g -L$(CURDIR)/capstone/build -lcapstone'

LIBCAPSTONE_OBJ := capstone/build/libcapstone.a

$(LIBCAPSTONE_OBJ):
	if [ ! -e capstone/Makefile ]; then \
		git submodule update --init --recursive; \
	fi
	cd capstone && \
		cmake -B build -DCMAKE_BUILD_TYPE=Release -DCAPSTONE_USE_ARCH_REGISTRATION=1 -DCAPSTONE_ARCHITECTURE_DEFAULT=1 -DCAPSTONE_BUILD_SHARED_LIBS=1 -DCAPSTONE_BUILD_CSTOOL=0 && \
		cmake --build build

.DEFAULT_GOAL := update
.PHONY: update
update:
	./genspec ./capstone/build
	./genconst ./capstone/bindings/python/capstone

.PHONY: gotest
gotest: $(LIBCAPSTONE_OBJ)
	@rm -f *.SPEC.test
	$(GO_CGO_CFLAGS) $(GO_CGO_LDFLAGS) go test -v .
