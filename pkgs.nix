{ flakeInputs ? import ./flake-fetch.nix
, system ? builtins.currentSystem
}:

let
  inherit (flakeInputs) nixpkgs gomod2nix npmlock2nix gitignore nix-eval-jobs;
in

import nixpkgs {
  inherit system;
  overlays = [
    (import "${gomod2nix}/overlay.nix")

    (final: prev: (import "${gitignore}" { inherit (final) lib; }))

    (final: prev: {
      npmlock2nix = import npmlock2nix { pkgs = final; };
    })

    (final: prev: {
      nix-eval-jobs = final.callPackage nix-eval-jobs { };
    })

    (final: prev:
      let
        inherit (prev) lib;
        dirNames = lib.attrNames (lib.filterAttrs (pkgDir: type: type == "directory" && builtins.pathExists (./packages + "/${pkgDir}/default.nix")) (builtins.readDir ./packages));
      in
      builtins.listToAttrs (map
        (pkgDir: {
          value = final.callPackage (./packages + "/${pkgDir}") { };
          name = pkgDir;
        })
        dirNames)
    )
  ];
}
