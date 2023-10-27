{ flakeInputs ? import ./flake-fetch.nix
, system ? builtins.currentSystem
}:

let
  inherit (flakeInputs) nixpkgs gomod2nix npmlock2nix gitignore;
in

import nixpkgs {
  inherit system;
  overlays = [
    (import "${gomod2nix}/overlay.nix")

    (final: prev: (import "${gitignore}" { inherit (final) lib; }))

    (final: prev: {
      npmlock2nix = import npmlock2nix { pkgs = final; };
    })
  ];
}
