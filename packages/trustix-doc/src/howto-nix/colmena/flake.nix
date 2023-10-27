{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";

    trustix = {
      url = "github:nix-community/trustix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };
  outputs = { nixpkgs, flake-utils, trustix, ... }:
    let
      hostName = "demo.trustix.dev";
    in

    # Provide colmena
    (flake-utils.lib.eachDefaultSystem
      (system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          devShells.default = pkgs.mkShell {
            packages = [ pkgs.colmena ];
          };
        })
    ) // {

      # Our actual deployment
      colmena = {
        meta = {
          nixpkgs = import nixpkgs {
            system = "x86_64-linux";
            # Note that the overlay has to be applied manually when using Colmena
            overlays = [ trustix.overlays.default ];
          };
        };

        # Import the Trustix NixOS modules on all machines
        defaults = { pkgs, ... }: {
          imports = [
            trustix.nixosModules.trustix
          ];
        };

        "${hostName}" = {

          # Main Trustix daemon configuration
          services.trustix = {
            enable = true;

            # Signers & publishers are separate concepts as the same key
            # could potentially be used to publish multiple logs under different Trustix subprotocols.
            #
            # If you don't know what this means: While Trustix is build for Nix first, the core can be used for other ecosystems too, not just Nix.
            signers.my-signer = {
              type = "ed25519";
              ed25519.private-key-path = "/var/lib/my-trustix-key";
            };

            publishers = [
              {
                signer = "my-signer"; # Use the key configured above
                protocol = "nix"; # This publisher is using the Nix subprotocol

                # An arbitrary (string -> string) attrset with metadata about this log
                meta = {
                  upstream = "https://cache.nixos.org";
                };

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
            # The logID we got earlier
            logID = "453016597475f45532e0a22a448ea7e0fb915e950d3c8930bfd23d962d73f9c1";
          };

          deployment = {
            targetHost = hostName;
            targetUser = "root";

            # We are using the Colmena secrets facility to upload the keys to the remote
            # without ending up world readable in the Nix store.
            keys = {
              my-trustix-key = {
                keyFile = ./secrets/log-priv;

                # This mode is too open for a real world deployment but we don't want to deal
                # with the complexities of secrets management here.
                permissions = "0644";
                destDir = "/var/lib";
              };
            };

          };
        };
      };
    };
}
