let
  pkgs = import ./nix;

in pkgs.mkShell {
  buildInputs = [
    pkgs.nixpkgs-fmt
    pkgs.niv
  ];
}
