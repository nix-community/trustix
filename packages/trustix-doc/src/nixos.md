# Trustix - Usage via Nix

The easiest way to use Trustix is via the NixOS modules, though even they require some manual preparation in terms of generating keys.

This document will guide you through the very basic NixOS setup required both by log clients and log publishers.

How to actually publish/subscribe are laid out in other documents.

## Requisites

- A git checkout of trustix
``` sh
$ git clone https://github.com/tweag/trustix.git
```

- [NixOS](https://nixos.org)

## Making modules available

Add the module import to your `configuration.nix` as such:
``` nix
{ pkgs, config, ...}:
{
  imports = [
    /path/to/trustix/git/checkout/nixos
  ];
}
```

This will add all relevant services to your system (but not enable them) and adds packages to the pkgs set via an overlay.
