let
  pkgs = import ./nix;

  pythonEnv = pkgs.python3.withPackages (ps: [ ps.grpcio ps.grpcio-tools ps.setuptools ]);

in
pkgs.mkShell {

  # Speed up compilation
  CGO_ENABLED = "0";

  buildInputs = [
    pkgs.hivemind # Process monitoring in development
    pkgs.reflex # File watcher utility
    pkgs.go
    pkgs.gocode
    pkgs.gore

    pkgs.nixpkgs-fmt
    pkgs.gomod2nix

    # pkgs.libgit2
    pkgs.pkgconfig

    pkgs.protobuf

    pkgs.systemfd # Socket activation testing

    # For development scripts
    pythonEnv
  ];

}
