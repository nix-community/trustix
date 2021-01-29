let

  overlay = self: super: {
    trustix = self.callPackage ../default.nix { };
  };

in [
  (import "${(builtins.fetchGit {
    url = "https://github.com/tweag/gomod2nix.git";
    rev = "929d740884811b388acc0f037efba7b5bc5745e8";
  })}/overlay.nix")
  (import "${(builtins.fetchGit {
    url = "https://github.com/nix-community/poetry2nix.git";
    rev = "31e3d16d65e65a11ae8e1aaeb9c2e7748144617f";
  })}/overlay.nix")
  overlay
]
