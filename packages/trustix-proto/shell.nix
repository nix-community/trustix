{ pkgs ? import ../../pkgs.nix { } }:

pkgs.mkShell {
  buildInputs = [
    pkgs.protoc-gen-go
    pkgs.protoc-gen-connect-go
    pkgs.protoc-gen-doc
  ];
}
