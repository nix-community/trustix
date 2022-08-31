{ pkgs ? import ../../../../pkgs.nix { } }:
let

  shellDrv = pkgs.npmlock2nix.shell {
    src = import ./src.nix { inherit (pkgs) lib; };
  };

in
shellDrv.overrideAttrs (old: {

  shellHook = ''
    if [[ "$(readlink -f node_modules)" == ${builtins.storeDir}* ]]; then
      rm -f node_modules
    fi
  '' + old.shellHook;

  buildInputs = old.buildInputs ++ [
    pkgs.hivemind
  ];
})
