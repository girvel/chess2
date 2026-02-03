APP_NAME := chess2
WIN_CC := x86_64-w64-mingw32-gcc
DIST_DIR := dist

.PHONY: all clean build-win pack release

all: release

clean:
	rm -rf $(DIST_DIR)
	rm -f $(APP_NAME).zip

build-win:
	@echo "Building for Windows..."
	@mkdir -p $(DIST_DIR)
	CGO_ENABLED=1 CC=$(WIN_CC) GOOS=windows GOARCH=amd64 \
	go build -ldflags "-s -w -H=windowsgui" \
	-o $(DIST_DIR)/$(APP_NAME).exe main.go

pack:
	@echo "Packing assets..."
	cp -r sprites $(DIST_DIR)/
	
	@echo "Zipping..."
	cd $(DIST_DIR) && zip -r ../$(APP_NAME).zip .

release: clean build-win pack
	@echo "Done! Created $(APP_NAME).zip"
