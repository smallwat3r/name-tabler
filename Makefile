APP := name-tabler
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo 1.0.0)

PLATFORMS := linux darwin windows
ARCHES := amd64 arm64

export CGO_ENABLED=0

.PHONY: clean build pack release

clean:
	rm -rf dist

build:
	@echo ">> Building $(APP) $(VERSION)"
	gox \
	  -os="$(PLATFORMS)" \
	  -arch="$(ARCHES)" \
	  -output="dist/$(APP)_{{.OS}}_{{.Arch}}/$(APP)" \
	  -ldflags="-s -w -X main.nameTablerVersion=$(VERSION) -buildid=" \
	  -gcflags=all=-trimpath=$(PWD)

pack: build
	@echo ">> Packaging artifacts"
	@cd dist && \
	for d in $(APP)_*; do \
	  if echo "$$d" | grep -qi windows; then \
	    echo ">> Zipping Windows artifact: $$d"; \
	    cp "$$d/$(APP)" "$$d/$(APP).exe" 2>/dev/null || true; \
	    (cd "$$d" && zip -9r "../$${d}.zip" "$(APP).exe" README* LICENSE* >/dev/null); \
	  else \
	    echo ">> Creating tar.gz for $$d"; \
	    tar czf "$${d}.tar.gz" -C "$$d" $(APP) README* LICENSE* 2>/dev/null || \
	    tar czf "$${d}.tar.gz" -C "$$d" $(APP); \
	  fi; \
	done

release: clean pack
	@echo "Artifacts in ./dist ready."
