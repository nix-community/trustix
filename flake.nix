{
  description = "Trustix";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    gomod2nix.url = "github:tweag/gomod2nix";

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

          devShells.default = pkgs.mkShell {
            packages = [
              (
                let
                  errorMessage = ''
                    Developing Trustix using Flakes is unsupported.

                    We are using the stable nix-shell interface together with direnv to recursively
                    load development shells for subpackages and relying on relative environment variables
                    for state directories and such, something which is not supported using Flakes.

                    For supported development methods see ./packages/trustix-doc/src/hacking.md.
                  '';
                in
                pkgs.runCommand "flakes-nein-danke" { } "echo '${errorMessage}' && exit 1"
              )
            ];
          };

        })
    );
}
