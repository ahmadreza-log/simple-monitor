APP_NAME   = simple-monitor
OUTPUT_DIR = build
OSES       = linux windows darwin
ARCHS      = amd64 arm64 386

.PHONY: clean build-all

clean:
	rm -rf $(OUTPUT_DIR)

build-all: clean
	mkdir -p $(OUTPUT_DIR)
	@for os in $(OSES); do \
	  for arch in $(ARCHS); do \
	    echo "Building for $$os/$$arch..."; \
	    ext=""; \
	    if [ $$os = "windows" ]; then ext=".exe"; fi; \
	    GOOS=$$os GOARCH=$$arch go build -o $(OUTPUT_DIR)/$(APP_NAME)-$$os-$$arch$$ext main.go; \
	  done \
	done
	@echo "✅ All builds are in '$(OUTPUT_DIR)'"