# Trustix - Using with Nix

The easiest way to use Trustix is via the NixOS modules.

This document walks you through how to subscribe to an already published binary cache.

## Configuring

- Add log(s) to your `configuration.nix`
``` nix
{ pkgs, config, ... }:
{

  services.trustix = {
    enable = true;

    subscribers = [
      {
        protocol = "nix";
        publicKey = {
          type = "ed25519";
          # Contents of the generated trustix-pub
          pub = "2uy8gNIOYEewTiV7iB7cUxBGpXxQtdlFepFoRvJTCJo=";
        };
      }
    ];

    # A remote can expose many logs and they are not neccesarily created by the remote in question
    remotes = [
      "grpc+https://trustix.example.com"
    ];

  };

}
```
