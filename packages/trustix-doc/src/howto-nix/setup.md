# Trustix - Usage via Nix

The easiest way to use Trustix is via the NixOS modules, though even they require some manual preparation in terms of generating keys.

This document will guide you through the very basic NixOS setup required both by log clients and log publishers.

How to actually publish/subscribe are laid out in other documents.

## Requisites
- A NixOS installation (flakes based optional)

## Strategies

### Classical Nix
It's highly recommended to use some automated tool like [niv](https://github.com/nmattia/niv) to ensure you are up to date with your external dependencies, here we'll show you how to integrate Trustix in your NixOS configuration _manually_ using no external tooling.

From within your configuration directory, clone Trustix:
``` sh
$ git clone https://github.com/nix-community/trustix.git
```

And add it to your NixOS configuration like:
```
{ config, pkgs, lib, ... }:
{
  imports = [ ./trustix/nixos ];
}
```

### Flakes
This is a minimal `flake.nix` for using Trustix with Flakes:
```

{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";

    trustix = {
      url = "github:nix-community/trustix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };
  outputs = { nixpkgs, flake-utils, trustix, ... }: {

    nixosConfigurations.trustix-example = nixpkgs.lib.nixosSystem {
      system = "x86_64-linux";
      modules =
        [ ({ pkgs, ... }: {
            # import trustix modules
            imports = [ trustix.nixosModules.trustix ];
          })
        ];
    };

  };
}
```

## Effect
This will add all relevant services to your system (but not enable them) and adds packages to the pkgs set via an overlay.
