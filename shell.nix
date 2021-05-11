{ pkgs ? import ./nix }:
let
  STATE_DIR = "${builtins.toString ./.}/state";
  TRUSTIX_RPC = "unix://${STATE_DIR}/trustix.sock";
  TRUSTIX_ROOT = builtins.toString ./.;

in
pkgs.mkShell {

  buildInputs = [
    pkgs.nixpkgs-fmt
    pkgs.hivemind
    pkgs.niv
  ];

  inherit STATE_DIR TRUSTIX_RPC TRUSTIX_ROOT;

}
