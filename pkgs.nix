{ flakeInputs ? (builtins.getFlake "${builtins.toString ./.}").inputs
, system ? builtins.currentSystem
}:

let
  inherit (flakeInputs) nixpkgs gomod2nix npmlock2nix;
in

import nixpkgs {
  inherit system;
  config.allowAliases = false;
  overlays = [
    gomod2nix.overlays.default

    (final: prev: {
      # Prevent the entirety of hydra to be in $PATH/runtime closure
      # We only want the evaluator
      hydra-eval-jobs = prev.runCommand "hydra-eval-jobs-${prev.hydra_unstable.version}" { } ''
        mkdir -p $out/bin
        cp -s ${prev.hydra_unstable}/bin/hydra-eval-jobs $out/bin/
      '';
    })

    (final: prev: {
      npmlock2nix = import npmlock2nix { pkgs = final; };
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
