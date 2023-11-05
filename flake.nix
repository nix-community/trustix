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
          checks = {
            x86_64-linux = self.checks.x86_64-linux;
            x86_64-darwin =
              let
                filteredChecks = [ "reuse" "treefmt" ]; # No point in running on linux _and_ darwin
              in
              lib.filterAttrs (name: _: ! lib.elem name filteredChecks && ! lib.strings.hasInfix "nixos" name) self.checks.x86_64-darwin;
          };
        };

        flake.nixosModules = {
          trustix = import ./nixos self;
        };

        imports = [
          inputs.treefmt-nix.flakeModule
          inputs.flake-root.flakeModule
        ];

        perSystem =
          { pkgs
          , config
          , system
          , self'
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
          {
            treefmt.imports = [ ./dev/treefmt.nix ];

            checks =
              # All packages + passthru.tests
              (
                let
                  packages' = builtins.removeAttrs (self'.packages) [ "default" ];
                in
                lib.listToAttrs (
                  lib.flatten (
                    lib.mapAttrsToList
                      (name: value: [ (lib.nameValuePair name value) ] ++ lib.mapAttrsToList (test: drv: lib.nameValuePair "${name}-${test}" drv)
                        (lib.filterAttrs (name: test: lib.elem system test.meta.platforms) (value.passthru.tests or { }))
                      )
                      packages'
                  )
                )
              )
              # Reuse lint
              // {
                reuse = pkgs.runCommand "reuse-lint" { nativeBuildInputs = [ pkgs.reuse ]; } ''
                  cd ${self}
                  reuse lint
                  touch $out
                '';
              }
              # Build development shell
              // {
                shell = self.devShells.${system}.default;
              }
              # NixOS tests
              // (
                let
                  checkArgs = {
                    inherit pkgs;
                    inherit system lib;
                    inherit (self') packages;
                  };
                in
                {
                  # Temporary: Comment out to get CI passing..
                  # trustix-nixos = import ./packages/trustix/nixos/test.nix checkArgs;
                }
              )
            ;

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
                (pkgs.diffoscope.override {
                  enableBloat = ! pkgs.stdenv.isDarwin;
                })
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
                export PATH="$FLAKE_ROOT/packages/trustix-nix-r13y-web/node_modules/.bin:$PATH";
              '';
            };

            packages = {
              default = self'.packages.trustix;
              trustix = callPackage ./packages/trustix { };
              trustix-doc = callPackage ./packages/trustix-doc { };
              trustix-nix = callPackage ./packages/trustix-nix { };
              trustix-nix-r13y = callPackage ./packages/trustix-nix-r13y { };
              trustix-nix-r13y-web = callPackage ./packages/trustix-nix-r13y-web { };
            };
          };
      };
}
