let
  sources = import ./sources.nix;
in
import sources.nixpkgs {
  overlays = [
    (import "${sources.naersk}/overlay.nix")
    (import "${sources.gomod2nix}/overlay.nix")
    (self: super: {
      npmlock2nix = import sources.npmlock2nix { pkgs = self; };
    })
    (self: super: {
      npmlock2nix = import sources.npmlock2nix { pkgs = self; };
    })

    (self: super: {
      # Prevent the entirety of hydra to be in $PATH/runtime closure
      # We only want the evaluator
      hydra-eval-jobs = self.runCommand "hydra-eval-jobs-${self.hydra-unstable.version}" { } ''
        mkdir -p $out/bin
        cp -s ${self.hydra-unstable}/bin/hydra-eval-jobs $out/bin/
      '';
    })

    # Local packages
    (self: super:
      let
        inherit (super) lib;
        dirNames = lib.attrNames (lib.filterAttrs (pkgDir: type: type == "directory" && builtins.pathExists (../packages + "/${pkgDir}/default.nix")) (builtins.readDir ../packages));
      in
      (
        builtins.listToAttrs (map
          (pkgDir: {
            value = self.callPackage (../packages + "/${pkgDir}") { };
            name = pkgDir;
          })
          dirNames)
      ))

  ];
}
