{
  description = "Trustix: Distributed trust and reproducibility tracking for binary caches";

  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "nixpkgs/3a7674c896847d18e598fa5da23d7426cb9be3d2";
    flake-compat = { url = "github:edolstra/flake-compat"; flake = false; };
    gomod2nix-flake = { url = "github:tweag/gomod2nix"; inputs.nixpkgs.follows = "nixpkgs"; };
    naersk-flake = { url = "github:nmattia/naersk"; inputs.nixpkgs.follows = "nixpkgs"; };
  };

  outputs =
    inputs@{ self
    , nixpkgs
    , flake-utils
    , flake-compat
    , gomod2nix-flake
    , naersk-flake
    }:
    { }
    //
    (flake-utils.lib.eachSystem [ "x86_64-linux" "x86_64-darwin" ]
      (system:
      let
        pkgs = import nixpkgs
          {
            inherit system;
            overlays = [
              self.overlay
              gomod2nix-flake.overlay
              naersk-flake.overlay
            ];
            config = {
              allowUnsupportedSystem = true; #enable for multi-arch check
              allowBroken = true;
            };
          };
      in
      rec {
        devShell = import ./shell.nix { inherit pkgs; };
        defaultPackage = pkgs.trustix;
        packages = {
          inherit (pkgs)
            nix-nar-unpack trustix trustix-nix trustix-nix-reprod;
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
          dirNames = lib.attrNames (lib.filterAttrs (pkgDir: type: type == "directory" && builtins.pathExists (./packages + "/${pkgDir}/default.nix")) (builtins.readDir ./packages));
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
