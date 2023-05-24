get:
	go get github.com/gotk3/gotk3
	go get github.com/gotk3/gotk3/gdk
	go get "github.com/sirupsen/logrus"

build:
	go build -v -o bin/nwg-look .

install:
	mkdir -p $(DESTDIR)/usr/share/nwg-look
	mkdir -p $(DESTDIR)/usr/share/nwg-look/langs
	mkdir -p $(DESTDIR)/usr/bin
	mkdir -p $(DESTDIR)/usr/share/applications
	mkdir -p $(DESTDIR)/usr/share/pixmaps
	cp stuff/main.glade $(DESTDIR)/usr/share/nwg-look/
	cp langs/* $(DESTDIR)/usr/share/nwg-look/langs/
	cp stuff/nwg-look.desktop $(DESTDIR)/usr/share/applications/
	cp stuff/nwg-look.svg $(DESTDIR)/usr/share/pixmaps/
	cp bin/nwg-look $(DESTDIR)/usr/bin

uninstall:
	rm -r $(DESTDIR)/usr/share/nwg-look
	rm $(DESTDIR)/usr/share/applications/nwg-look.desktop
	rm $(DESTDIR)/usr/share/pixmaps/nwg-look.svg
	rm $(DESTDIR)/usr/bin/nwg-look

run:
	go run .
