let
  pkgs = import ./pkgs.nix { };
  inherit (pkgs) lib;

  STATE_DIR = "${builtins.toString ./.}/state";
  TRUSTIX_RPC = "unix://${STATE_DIR}/trustix.sock";
  TRUSTIX_ROOT = builtins.toString ./.;

  # Wrap treefmt with a Go compiler so it can do gofmt without recursively loading subprojects
  treefmt = pkgs.writeShellScriptBin "treefmt" ''
    export PATH=${pkgs.go}/bin:$PATH
    exec ${pkgs.treefmt}/bin/treefmt "$@"
  '';

  python = pkgs.python3.override {
    self = python;
    packageOverrides = self: super: {

      pretty-errors =
        let
          version = "1.2.25";
        in
        self.buildPythonPackage {
          pname = "pretty-errors";
          inherit version;

          src = self.fetchPypi {
            pname = "pretty_errors";
            inherit version;
            hash = "sha256-oWulx1LIfCY7+S+LS1hiTjseKScak5H1ZPErhuk8Z1U=";
          };

          # Work around interactive installer
          postPatch = "rm ./pretty_errors/__main__.py";

          propagatedBuildInputs = [
            self.colorama
          ];
        };

    };
  };

  # Some development tools (like the license file generator) is written in Python
  pythonEnv = python.withPackages (ps: [
    ps.mypy
    ps.black
    ps.pretty-errors
  ]);

  sqlFormatterWriter = pkgs.writeScriptBin "sql-formatter-writer" ''
    #!${pkgs.runtimeShell}
    set -euo pipefail
    exec ${pkgs.nodePackages.sql-formatter}/bin/sql-formatter -l sqlite "$1" | ${pkgs.moreutils}/bin/sponge "$1"
  '';

in
pkgs.mkShell {

  buildInputs = [
    # Development scripts
    pythonEnv

    # Meta code formatter
    treefmt

    # Only build job if it's not in the binary cache
    pkgs.nix-build-uncached

    # Protobuf formatter (clang-format)
    pkgs.clang-tools

    # Format Nix expressions
    pkgs.nixpkgs-fmt

    # Format SQL
    sqlFormatterWriter

    # Procfile process runner
    pkgs.hivemind

    # Nix go modules code generator
    pkgs.gomod2nix

    # Protobuf
    pkgs.protobuf
    pkgs.grpcurl # gRPC CLI

    # Go linters
    pkgs.golangci-lint # Multi purpose linter

    # File system watchers
    pkgs.reflex
    pkgs.entr

    # Docs
    pkgs.mdbook

    # License management and compliance
    pkgs.reuse
  ];

  inherit STATE_DIR TRUSTIX_RPC TRUSTIX_ROOT;

  # Write token used for log submission
  TRUSTIX_TOKEN = "${builtins.toString ./packages/trustix/dev/token-priv}";

}
