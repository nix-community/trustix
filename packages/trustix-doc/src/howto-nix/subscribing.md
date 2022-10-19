# Trustix - Subscribing
This document walks you through how to subscribe to an already published binary cache.

## Requisites
- A local Trustix instance
- A remote log's metadata
  - Public key
  - URL

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
          key = "2uy8gNIOYEewTiV7iB7cUxBGpXxQtdlFepFoRvJTCJo=";
        };
      }
    ];

    # A remote can expose many logs and they are not neccesarily created by the remote in question
    remotes = [
      "https://demo.trustix.dev"
    ];

  };

}
```
