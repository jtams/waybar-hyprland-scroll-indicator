PREFIX ?= ~/.local

install:
	go build -o $(PREFIX)/bin/hypr-scroll-indicator .
	@echo "Installed hypr-scroll-indicator to $(PREFIX)/bin/hypr-scroll-indicator"
