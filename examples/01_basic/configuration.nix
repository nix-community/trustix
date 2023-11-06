{ config, pkgs, lib, ... }:
{
  services.trustix = {
    enable = true;

    # Trustix differentiates between the concepts of a "signer" and a "publisher".
    # A signer refers to a private key implementation.
    # These can be file based or use hardware tokens.
    signers.my-signer = {
      type = "ed25519";
      # Configuring the private key like this by path is bad practice because the key ends up world readable in /nix/store.
      # You should either:
      # - Put the key in a persistent path and reference it like: `ed25519.private-key-path = "/path/to/key"`
      # - Use a secrets management solution like sops-nix or agenix.
      ed25519.private-key-path = ./secrets/log-priv;
    };

    publishers = [
      {
        # Use the key configured above
        signer = "my-signer";

        # Trustix is built first and foremost for Nix, but could also be used for verifying other package ecosystems.
        protocol = "nix";

        # An arbitrary (string -> string) attrset with metadata about this log.
        # This isn't used by the Trustix logs but is used to inform the Nix binary cache proxy about possible substitution sources.
        meta = {
          upstream = "https://cache.nixos.org";
        };

        # The public key identifying this log.
        publicKey = {
          type = "ed25519";
          key = builtins.readFile ./secrets/log-pub;
        };
      }
    ];
  };
}
