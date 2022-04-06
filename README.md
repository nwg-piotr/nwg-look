# nwg-look
This application is a part of the [nwg-shell](https://github.com/nwg-piotr/nwg-shell) project.

Nwg-look is a GTK3 settings editor, designed to work properly in wlroots-based Wayland environment. The look and feel is strongly influenced by [LXAppearance](https://wiki.lxde.org/en/LXAppearance), but nwg-look is intended to free the user from a few inconveniences:

- It works natively on Wayland. You no longer need Xwayland, nor strange env variables for it to run.
- It applies gsettings directly, with no need to use [workarounds](https://github.com/swaywm/sway/wiki/GTK-3-settings-on-Wayland). You don't need to set gsettings in the sway config file. You don't need the `import-gsettings` script.

![screenshot](https://user-images.githubusercontent.com/20579136/161869170-ef1abcfd-c72c-4da9-8cee-1f9560d2b5af.png)

## Dependencies

- go (just to build)
- gtk3
- [xcur2png](https://github.com/eworm-de/xcur2png)

## Installation

1. Clone the repository, cd into it.
2. `make build`
3. `sudo make install`

## Usage

```text
$ nwg-look -h
Usage of nwg-look:
  -a	Apply stored gsettings and quit
  -d	turn on Debug messages
  -n	do Not save gtk settings.ini
  -v	display Version information
```

The `-a` flag has been added just in case. When you press the "Apply" button, in addition to applying the changes, a backup file is also created. You may apply gsetting again w/o running the GUI, by just `nwg-look -a`. No idea if it's going to be useful in real life. ;)

## Development status

This is the very first public release. Bugs and missing features are expected. Thanks in advance for reporting them.
