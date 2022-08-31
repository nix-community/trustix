{
  description = "Trustix";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    gomod2nix.url = "github:tweag/gomod2nix";

    npmlock2nix = {
      url = "github:nix-community/npmlock2nix/master";
      flake = false;
    };
  };

  outputs = { self, nixpkgs, flake-utils, gomod2nix, npmlock2nix }@flakeInputs:
    (flake-utils.lib.eachDefaultSystem
      (system:
        let
          pkgs = import ./pkgs.nix {
            inherit system flakeInputs;

          };
        in
        {
          packages = {
            inherit (pkgs) trustix trustix-doc trustix-nix trustix-nix-reprod;
          };
        })
    );
}
