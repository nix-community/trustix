# Trustix - Usage via Nix

The easiest way to use Trustix is via the NixOS modules, though even they require some manual preparation in terms of generating keys.

This document will guide you through the very basic NixOS setup required both by log clients and log publishers.

How to actually publish/subscribe are laid out in other documents.

## Requisites
- A NixOS installation using Flakes

## Create keys

All Trustix build logs are first and foremost identified by their key pair, which will be the first thing we have to generate.

Let's start by generating a key pair for our log:
```
$ mkdir secrets
$ nix run github:nix-community/trustix#trustix -- generate-key --privkey secrets/log-priv --pubkey secrets/log-pub
```

Additionally logs are identified not just by their key, but how that key is used.
If a key is used for multiple protocols (not just Nix) those logs will have a different ID.
This ID is what _subscribers_ use to indicate what they want to subscribe to.

To find out the log ID for the key pair you just generated:
`$ nix run github:nix-community/trustix#trustix -- print-log-id --protocol nix --pubkey $(cat secrets/log-pub)`

## Flakes

- `flake.nix`
``` nix
{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

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
            imports = [
              trustix.nixosModules.trustix
              ./configuration.nix
            ];
          })
        ];
    };

  };
}
```

- `configuration.nix`:
``` nix
{{#include ../../../../examples/01_basic/configuration.nix}}
```

## Effect
This will set up an instance of Trustix on your system.
In the next chapter we will look at using the post build hook to publish results to our local log.
