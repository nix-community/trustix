{ pkgs ? import ../../pkgs.nix { } }:

pkgs.mkShell {
  buildInputs = [
    pkgs.hivemind
    pkgs.protobuf
    pkgs.reflex
    pkgs.go

    pkgs.protoc-gen-doc
    pkgs.mdbook
  ];

  shellHook = ''
    export GOBIN=$(mktemp -d)
    export PATH=$GOBIN:$PATH
    go list -f '{{range .Imports}}{{.}} {{end}}' tools.go | xargs go install
  '';
}
