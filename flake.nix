{
  description = "Trustix";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    gomod2nix.url = "github:nix-community/gomod2nix";

    npmlock2nix = {
      url = "github:nix-community/npmlock2nix/master";
      flake = false;
    };
  };

  outputs = { self, nixpkgs, flake-utils, gomod2nix, npmlock2nix }@flakeInputs:
    (flake-utils.lib.eachDefaultSystem
      (system:
        let
          pkgs = import ./pkgs.nix {
            inherit system flakeInputs;

          };
        in
        {
          packages = import ./default.nix { inherit pkgs; };

          # Fake shell derivation that evaluates but doesn't build and producec an error message
          # explaining the supported setup.
          devShells.default =
            let
              errorMessage = ''
                Developing Trustix using Flakes is unsupported.

                We are using the stable nix-shell interface together with direnv to recursively
                load development shells for subpackages and relying on relative environment variables
                for state directories and such, something which is not supported using Flakes.

                For supported development methods see ./packages/trustix-doc/src/hacking.md.
              '';
            in
            builtins.derivation {
              name = "flakes-nein-danke-shell";
              builder = "bash";
              inherit system;
              preferLocalBuild = true;
              allowSubstitutes = false;
              fail = builtins.derivation {
                name = "flakes-nein-danke";
                builder = "/bin/sh";
                args = [ "-c" "echo '${errorMessage}' && exit 1" ];
                preferLocalBuild = true;
                allowSubstitutes = false;
                inherit system;
              };
            };
        }
      )
    );
}
