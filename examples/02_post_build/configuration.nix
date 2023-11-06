{ config, pkgs, lib, ... }:
{
  # Our basic Trustix configuration from before
  services.trustix = {
    enable = true;

    signers.my-signer = {
      type = "ed25519";
      ed25519.private-key-path = ./secrets/log-priv;
    };

    publishers = [
      {
        signer = "my-signer";
        protocol = "nix";
        meta.upstream = "https://cache.nixos.org";
        publicKey = {
          type = "ed25519";
          key = builtins.readFile ./secrets/log-pub;
        };
      }
    ];
  };

  # Enable the post build hook to push builds to the main Trustix daemon
  services.trustix-nix-build-hook = {
    enable = true;
    # Log id as returned by `trustix print-log-id --protocol nix --pubkey $(cat secrets/log-pub)`
    # This is your logs globally unique identifier and what clients will use to subscribe to your build results.
    logID = "0c7942343fa91b610704d531f552f3e785705dbd7d22c965bc0d58fa3ff2c87c";
  };
}
