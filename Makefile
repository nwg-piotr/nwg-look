get:
	go get github.com/gotk3/gotk3
	go get github.com/gotk3/gotk3/gdk
	go get "github.com/sirupsen/logrus"

build:
	go build -o nwg-look .

install:
	mkdir -p /usr/share/nwg-look
	cp stuff/main.glade /usr/share/nwg-look/
	cp stuff/nwg-look.desktop /usr/share/applications/
	cp stuff/nwg-look.svg /usr/share/pixmaps/
	cp nwg-look /usr/bin

uninstall:
	rm -r /usr/share/nwg-look
	rm /usr/share/applications/nwg-look.desktop
	rm /usr/share/pixmaps/nwg-look.svg
	rm /usr/bin/nwg-look

run:
	go run .
