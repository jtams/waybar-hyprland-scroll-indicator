This adds a scroll indicator for hyprland's new `scrolling` layout.

I personally prefer to keep the scrolling layout columns 100% wide, but that makes it impossible to know if there's any window to the left or right. This plugin fixes that.

# Installation

## Requirements

- Waybar v0.15.0+
- Hyprland v0.54.0+
- Go v1.26.0+

These are just the versions I use. It almost certainly works on older versions.

## Installation

```
make install
```

This will install it to ~/.local/bin/hypr-scroll-indicator

You can change the path in the Makefile.

## Usage

In your Waybar config add:

```
"custom/scroll_indicator": {
    "exec": "~/.local/bin/hypr-scroll-indicator",
    "return-type": "json",
    "format": "{}"
},
```

And add it to the module list of your choice:

```
"modules-center": ["custom/scroll_indicator"],
```

Restart Waybar.

# License

MIT
