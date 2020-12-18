let
  src =
    let
      meta = builtins.fromJSON (builtins.readFile ./nixpkgs.json);
    in
    builtins.fetchTarball {
      url = "https://github.com/nixos/nixpkgs-channels/archive/${meta.rev}.tar.gz";
      sha256 = meta.sha256;
    };

  args = {
    overlays = [
      (import ./overlay.nix)
      (import "${(builtins.fetchGit {
        url = "git@github.com:tweag/gomod2nix.git";
        rev = "929d740884811b388acc0f037efba7b5bc5745e8";
      })}/overlay.nix")
      (import "${(builtins.fetchGit {
        url = "git@github.com:nix-community/poetry2nix.git";
        rev = "31e3d16d65e65a11ae8e1aaeb9c2e7748144617f";
      })}/overlay.nix")
    ];
  };

  patches = [ ];

  pkgs = import src args;

  patched = import
    (pkgs.stdenv.mkDerivation {
      name = "nixpkgs";
      inherit src patches;
      dontBuild = true;
      preferLocalBuild = true;
      fixupPhase = ":"; # We dont need to patch nixpkgs shebangs
      installPhase = ''
        mkdir -p $out
        cp -a .version * $out/
      '';
    })
    args;

in
if patches == [ ] then pkgs else patched
