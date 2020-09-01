let
  pkgs = import ./nix;

  pythonEnv = pkgs.python3.withPackages(ps: []);

in pkgs.mkShell {

  buildInputs = [
    pkgs.hivemind  # Process monitoring in development
    pkgs.reflex  # File watcher utility
    pkgs.go

    pkgs.libgit2
    pkgs.pkgconfig

    pkgs.protobuf

    # For development scripts
    pythonEnv

  ];

  shellHook = ''
    unset GOPATH
  '';

}
