# Trustix - Binary cache setup

The easiest way to use Trustix is via the NixOS modules, though even they require some manual preparation in terms of generating keys.

This document walks you through how to configure your local system as a binary cache.

## Requisites

We are assuming you have already followed the steps to set up one or more subscribers to your local Trustix instance.

- Generate a public/private keypair to use with your local binary cache.
``` sh
$ nix-store --generate-binary-cache-key binarycache.example.com cache-priv-key.pem cache-pub-key.pem
```

- Move the keys somewhere persistent and safe
Of course having keys around readable by anyone on the system is not a good idea, so we will move these somewhere safe.
In this tutorial we are using `/keys` but you are free to use whatever you wish.

`$ mv cache-priv-key.pem /keys/cache-priv-key.pem`

## Configuring

- Add the binary cache to your `configuration.nix`
``` nix
{ pkgs, config, ... }:
{

  # Enable the local binary cache server
  services.trustix-nix-cache = {
    enable = true;
    private-key = "/keys/cache-priv-key.pem";
    port = 9001;
  };

  # Configure Nix to use it
  nix = {
    binaryCaches = [
      "http//localhost:9001"
    ];
    binaryCachePublicKeys = [
      "binarycache.example.com://06YZJreoL8n9IdDlhnA3t7uJmHUI/rIIy3uO4FHRY="
    ];
  };

  # Configure your Trustix daemon with a decision making process on how
  # to determine if a build is trustworthy or not.
  #
  # In this case we configure it to have at least 2/3 majority to be substituted.
  #
  # Note that this configuration is incomplete and assumes you have already set up a subscriber.
  services.trustix = {
    deciders.nix = [
      {
        {
          type = "percentage";
          percentage.minimum = 66;
        }
      }
    ];
  };

}
```

You are now all set up to use Trustix as a substitution method!
