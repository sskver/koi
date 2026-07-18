# Per-platform build settings for everything under cmd/<platform>/.
# Add a new board by adding its directory under cmd/ plus a line here.
GOOS_rpi   := linux
GOARCH_rpi := arm64

PLATFORMS := $(notdir $(wildcard cmd/*))
DEMOS     := $(foreach p,$(PLATFORMS),$(notdir $(wildcard cmd/$(p)/*)))

.PHONY: test build clean $(DEMOS)

test:
	go test ./...

build: $(DEMOS)

bin:
	mkdir -p bin

define DEMO_RULE
$(2): bin
	CGO_ENABLED=0 GOOS=$$(GOOS_$(1)) GOARCH=$$(GOARCH_$(1)) go build -o bin/$(2) ./cmd/$(1)/$(2)
endef
$(foreach p,$(PLATFORMS),$(foreach d,$(notdir $(wildcard cmd/$(p)/*)),$(eval $(call DEMO_RULE,$(p),$(d)))))

clean:
	rm -rf bin
