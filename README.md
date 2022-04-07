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
  -a		Apply stored gsetting and quit
  -d		turn on Debug messages
  -n		do Not save gtk settings.ini
  -r		Restore default values and quit
  -v		display Version information
```

The `-a` flag has been added just in case. When you press the "Apply" button, in addition to applying the changes, a backup file is also created. You may apply gsetting again w/o running the GUI, by just `nwg-look -a`. No idea if it's going to be useful in real life. ;)

### Usage in sway

The default way to apply GTK setting on [sway](https://github.com/swaywm/sway) Wayland compositor has been
described in the [GTK 3 settings on Wayland](https://github.com/swaywm/sway/wiki/GTK-3-settings-on-Wayland)
Wiki section. **You no longer need it**. Nwg-look loads and saves gsettings values directly, and does not
care about the `~/.config/gtk-3.0/settings.ini` values. It only exports your settings to it, unless you use
the `-n` flag.

Therefore, if your sway config file contains either

```text
set $gnome-schema org.gnome.desktop.interface

exec_always {
    gsettings set $gnome-schema gtk-theme 'Your theme'
    gsettings set $gnome-schema icon-theme 'Your icon theme'
    gsettings set $gnome-schema cursor-theme 'Your cursor Theme'
    gsettings set $gnome-schema font-name 'Your font name'
}
```

or if you use the `import-gsettings` script:

```text
exec_always import-gsettings
```

to parse and apply the settings.ini file, **remove these lines**.

## Backward compatibility

Some gsetting keys have no direct counterparts in the Gtk.Settings type. While exporting
the settings.ini file, nwg-look uses the most similar values:

| gsettings | Gtk.Settings |
| --------- | ------------ |
| `font-hinting` | `gtk-xft-hintstyle` |
| `font-antialiasing` | `gtk-xft-antialias` |
| `font-rgba-order` | `gtk-xft-rgba` |

Some **Other** settings have been left just for LXAppearance compatibility, and possible
use of your settings.ini file elsewhere:

- Toolbar style
- Toolbar icon size

have been deprecated since GTK 3.10, and the values are ignored.

- Show button images
- Show menu images

have been deprecated since GTK 3.10, and have no corresponding gsettings values.

- Enable event sounds
- Enable input feedback sounds

don't seem to change anything in non-GNOME environment.

## Development status

This is the very first public release. Bugs and missing features are expected. Thanks in advance for reporting them.
