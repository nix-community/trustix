{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {

  buildInputs = [
    pkgs.go-ethereum
    pkgs.hivemind  # Process monitoring in development
    pkgs.solc  # Solidity compiler
    pkgs.reflex  # File watcher utility
    pkgs.go
    pkgs.vgo2nix
  ];

  shellHook = ''
    unset GOPATH
  '';

}
