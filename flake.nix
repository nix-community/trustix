{
  description = "A very basic flake";

  inputs = {
    # TODO: need to pin this to the source rev
    nixpkgs.url =
      "github:NixOS/nixpkgs/f5e8bdd07d1afaabf6b37afc5497b1e498b8046f";
    # nixpkgs.url = "github:NixOS/nixpkgs/release-21.05";
    npmlock2nix.url =
      "github:nix-community/npmlock2nix/7a321e2477d1f97167847086400a7a4d75b8faf8";
    gomod2nix.url =
      "github:tweag/gomod2nix/c78d7b9f15a24eba95fbc228509f513c83709d8b";
    flake-utils.url = "github:numtide/flake-utils";
    npmlock2nix.flake = false;
  };

  outputs = inputs@{ self, nixpkgs, npmlock2nix, gomod2nix, flake-utils }:
    (flake-utils.lib.eachDefaultSystem (system:
      let
        lib = pkgs.lib;
        pkgs = import nixpkgs {
          inherit system;

          # TODO: this was not working with overlay.default
          # error said  "it is a function expected a set"
          overlays = [
            gomod2nix.overlay
            (self: super: {
              npmlock2nix = import npmlock2nix { pkgs = self; };
            })
            (self: super: {
              # Prevent the entirety of hydra to be in $PATH/runtime closure
              # We only want the evaluator
              hydra-eval-jobs =
                self.runCommand "hydra-eval-jobs-${self.hydra-unstable.version}"
                { } ''
                  mkdir -p $out/bin
                  cp -s ${self.hydra-unstable}/bin/hydra-eval-jobs $out/bin/
                '';
            })
          ];
        };
      in {
        packages.default = pkgs.callPackage ./packages/trustix-doc { };
        packages.trustix = pkgs.callPackage ./packages/trustix { };
        packages.trustix-nix = pkgs.callPackage ./packages/trustix-nix { };
        packages.trustix-doc = pkgs.callPackage ./packages/trustix-doc { };
        packages.trustix-nix-reprod =
          pkgs.callPackage ./packages/trustix-nix-reprod { };
      }));
}
