let
  pkgs = import <nixpkgs> { overlays = import ../../../../nix/overlays.nix; };
  inherit (pkgs) lib;

in
pkgs.npmlock2nix.build {
  src = lib.cleanSourceWith {
    filter = name: type: !(lib.hasSuffix "node_modules" name) && !(lib.hasSuffix ".direnv" name) && !(lib.hasSuffix "dist" name);
    src = lib.cleanSource ./.;
  };

  buildCommands = [ "npm run build" ];
  installPhase = "cp -r dist $out";

}
