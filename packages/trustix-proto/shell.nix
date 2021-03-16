{ pkgs ? import ../../nix }:

pkgs.mkShell {
  buildInputs = [
    pkgs.hivemind
    pkgs.protobuf
    pkgs.reflex
    pkgs.go
  ];

  shellHook = ''
    export GOBIN=$(mktemp -d)
    export PATH=$GOBIN:$PATH
    go list -f '{{range .Imports}}{{.}} {{end}}' tools.go | xargs go install
  '';
}
