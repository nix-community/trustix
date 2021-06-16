{
  description = "Trustix: Distributed trust and reproducibility tracking for binary caches";

  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "nixpkgs/release-21.05";
    devshell-flake.url = "github:numtide/devshell";
    flake-compat = { url = "github:edolstra/flake-compat"; flake = false; };
    gomod2nix-flake = { url = "github:tweag/gomod2nix"; };
    naersk-flake = { url = "github:nmattia/naersk"; };
  };

  outputs =
    inputs@{ self
    , nixpkgs
    , devshell-flake
    , flake-utils
    , flake-compat
    , gomod2nix-flake
    , naersk-flake
    }:
    {
      nixosModules = {
        trustix = import ./packages/trustix/nixos;
        trustix-nix = import ./packages/trustix-nix/nixos;
      };
    }
    //
    (flake-utils.lib.eachDefaultSystem
      (system:
      let
        pkgs = import nixpkgs
          {
            inherit system;
            overlays = [
              self.overlay
              gomod2nix-flake.overlay
              naersk-flake.overlay
              devshell-flake.overlay
            ];
            config = {
              allowUnsupportedSystem = true; #enable for multi-arch check
              allowBroken = true;
            };
          };
      in
      rec {
        devShell = with pkgs;
          devshell.mkShell {
            imports = [
              (devshell.importTOML ./devshell.toml)
            ];
          };
        defaultPackage = pkgs.trustix;
        packages = {
          inherit (pkgs)
            trustix trustix-doc trustix-nix trustix-nix-reprod;
        };
        hydraJobs = {
          inherit packages;
        };
      }
      )
    ) //
    {
      overlay = final: prev:
        let
          inherit (prev) lib;
          dirNames = builtins.attrNames (builtins.readDir ./packages);
        in
        (
          builtins.listToAttrs (map
            (pkgDir: {
              value = prev.callPackage (./packages + "/${pkgDir}") { };
              name = pkgDir;
            })
            dirNames)
        ) // {

          hydra-eval-jobs = prev.runCommand "hydra-eval-jobs-${prev.hydra-unstable.version}" { } ''
            mkdir -p $out/bin
            cp -s ${prev.hydra-unstable}/bin/hydra-eval-jobs $out/bin/
          '';
        };
    };
}
