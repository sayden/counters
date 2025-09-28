.PHONY: all
all: docs
	wails build
	mv build/bin/counters-desktop ~/.local/bin
	cp ./build/counters-desktop.desktop ~/.local/share/applications
	cp ./build/appicon7.png ~/.local/share/icons/hicolor/1024x1024/apps/counters-desktop.png

.PHONY: docs
docs:
	cd documentation && bun run build
