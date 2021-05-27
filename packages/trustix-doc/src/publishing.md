# Trustix - Publishing with Nix

The easiest way to use Trustix is via the NixOS modules, though even they require some manual preparation in terms of generating keys.

This document walks you through creating key pairs and publishing a log.

## Requisites

The commands in this document assumes you have a git checkout of Trustix and that you already have the NixOS module import set up.

## Preparation

- Enter the Trustix git repository
``` sh
$ cd trustix
```

- Build the Trustix command to generate key material
`$ nix-build ./. -A packages.trustix`

- Generate an ed25519 keypair
`$ ./result/bin/trustix generate-key --privkey ./trustix-priv --pubkey ./trustix-pub`

This will create two files, `pub` and `priv`.

- Move the keys somewhere persistent and safe
Of course having keys around readable by anyone on the system is not a good idea, so we will move these somewhere safe.
In this tutorial we are using `/keys` but you are free to use whatever you wish.

`$ mv trustix-priv /keys/trustix-priv`

## Configuring

- Add log to your `configuration.nix`
``` nix
{ pkgs, config, ... }:
{

  services.trustix = {
    enable = true;

    signers.snakeoil = {
      type = "ed25519";
      ed25519 = {
        private-key-path = "/keys/trustix-priv";
      };
    };

    publishers = [
      {
        signer = "snakeoil";
        protocol = "nix";
        publicKey = {
          type = "ed25519";
          # Contents of the generated trustix-pub
          pub = "2uy8gNIOYEewTiV7iB7cUxBGpXxQtdlFepFoRvJTCJo=";
        };
      }
    ];

  };

  # Push local builds via the post-build hook
  services.trustix-nix-build-hook.enable = true;

}
```

You can now use Nix as normal and it will publish the locally performed builds to your log!
