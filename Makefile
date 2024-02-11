# For alternate install dir (e.g. "/usr/local")
# specify PREFIX in make command:
#   sudo make PREFIX=/usr/local install
# Defaults to "/usr" if not specified.
PREFIX ?= /usr

get:
	go get github.com/gotk3/gotk3
	go get github.com/gotk3/gotk3/gdk
	go get "github.com/sirupsen/logrus"

build:
	go build -v -o bin/nwg-look .

install:
	mkdir -p $(DESTDIR)$(PREFIX)/share/nwg-look
	mkdir -p $(DESTDIR)$(PREFIX)/share/nwg-look/langs
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	mkdir -p $(DESTDIR)$(PREFIX)/share/applications
	mkdir -p $(DESTDIR)$(PREFIX)/share/pixmaps

	mkdir -p $(DESTDIR)$(PREFIX)/share/doc/nwg-look
	mkdir -p $(DESTDIR)$(PREFIX)/share/licenses/nwg-look

	cp stuff/main.glade $(DESTDIR)$(PREFIX)/share/nwg-look/
	cp langs/* $(DESTDIR)$(PREFIX)/share/nwg-look/langs/
	cp stuff/nwg-look.desktop $(DESTDIR)$(PREFIX)/share/applications/
	cp stuff/nwg-look.svg $(DESTDIR)$(PREFIX)/share/pixmaps/
	cp bin/nwg-look $(DESTDIR)$(PREFIX)/bin

	cp README.md $(DESTDIR)$(PREFIX)/share/doc/nwg-look
	cp LICENSE $(DESTDIR)$(PREFIX)/share/licenses/nwg-look

uninstall:
	rm -r $(DESTDIR)$(PREFIX)/share/nwg-look
	rm $(DESTDIR)$(PREFIX)/share/applications/nwg-look.desktop
	rm $(DESTDIR)$(PREFIX)/share/pixmaps/nwg-look.svg
	rm $(DESTDIR)$(PREFIX)/bin/nwg-look

run:
	go run .
