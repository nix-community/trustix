let
  pkgs = import <nixpkgs> { overlays = import ../../nix/overlays.nix; };
  inherit (pkgs) poetry2nix;

in poetry2nix.mkPoetryApplication {
  projectDir = ./.;
  overrides = poetry2nix.overrides.withDefaults (
    import ./overrides.nix { inherit pkgs; }
  );
}
