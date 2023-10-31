{
  description = "Trustix";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };

    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    gitignore = {
      url = "github:hercules-ci/gitignore.nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    npmlock2nix = {
      url = "github:nix-community/npmlock2nix/master";
      flake = false;
    };

    nix-github-actions = {
      url = "github:nix-community/nix-github-actions";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    flake-root.url = "github:srid/flake-root";
  };

  outputs =
    { self
    , nixpkgs
    , flake-parts
    , gomod2nix
    , npmlock2nix
    , gitignore
    , systems
    , treefmt-nix
    , flake-root
    , nix-github-actions
    ,
    } @ inputs:
    let
      inherit (nixpkgs) lib;
    in
    flake-parts.lib.mkFlake
      { inherit inputs; }
      {
        systems = [
          "x86_64-linux"
          "aarch64-linux"
          "x86_64-darwin"
          "aarch64-darwin"
        ];

        flake.githubActions = nix-github-actions.lib.mkGithubMatrix {
          checks = { inherit (self.checks) x86_64-linux; };
        };

        flake.nixosModules = {
          trustix = import ./nixos self;
        };

        flake.overlays.default = final: prev: import ./default.nix { };

        imports = [
          inputs.treefmt-nix.flakeModule
          inputs.flake-root.flakeModule
        ];

        perSystem =
          { pkgs
          , config
          , system
          , ...
          }:
          let
            callPackage = lib.callPackageWith (pkgs
              // {
              inherit (inputs.gomod2nix.legacyPackages.${system}) buildGoApplication;
              inherit (inputs.gitignore.lib) gitignoreSource;
              npmlock2nix = import npmlock2nix { inherit pkgs; };
            });
          in
          rec {
            treefmt.imports = [ ./dev/treefmt.nix ];

            checks =
              (
                let
                  packages' = builtins.removeAttrs packages [ "default" ];
                in
                lib.listToAttrs (
                  lib.flatten (
                    lib.mapAttrsToList
                      (name: value: [ (lib.nameValuePair name value) ] ++ lib.mapAttrsToList (test: drv: lib.nameValuePair "${name}-${test}" drv)
                        (value.passthru.tests or { })
                      )
                      packages'
                  )
                )
              )
              // {
                reuse = pkgs.runCommand "reuse-lint" { nativeBuildInputs = [ pkgs.reuse ]; } ''
                  cd ${self}
                  reuse lint
                  touch $out
                '';
              }
              // {
                shell = self.devShells.${system}.default;
              };

            devShells.default = pkgs.mkShell {
              buildInputs = [
                # Procfile process runner
                pkgs.hivemind

                # Nix go modules code generator
                inputs.gomod2nix.packages.${system}.default

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

                # Socket activation testing
                pkgs.systemfd

                # Dev
                pkgs.go
                pkgs.nix-eval-jobs
                pkgs.sqlite
                pkgs.diffoscope
                pkgs.sqlc
                pkgs.goose
                pkgs.protoc-gen-go
                pkgs.protoc-gen-doc
                pkgs.protoc-gen-connect-go
                pkgs.nodejs
              ];

              inputsFrom = [ config.flake-root.devShell ];

              # Write token used for log submission
              env.TRUSTIX_TOKEN = "${./packages/trustix/dev/token-priv}";

              shellHook = ''
                export TRUSTIX_RPC="unix://$FLAKE_ROOT/state/trustix.sock"
                export TRUSTIX_NIX_REPROD_STATE_DIR="$FLAKE_ROOT/state/nix-reprod"
                export PATH=${builtins.toString ./packages/trustix-nix-r13y-web}/node_modules/.bin:$PATH
                export TRUSTIX_STATE_DIR="$FLAKE_ROOT/state/trustix";
              '';
            };

            packages = {
              trustix = callPackage ./packages/trustix { };
              trustix-doc = callPackage ./packages/trustix-doc { };
              trustix-nix = callPackage ./packages/trustix-nix { };
              trustix-nix-r13y = callPackage ./packages/trustix-nix-r13y { };
              trustix-nix-r13y-web = callPackage ./packages/trustix-nix-r13y-web { };
            };
          };
      };
}
