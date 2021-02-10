let

  overlay = self: super: {
    trustix = self.callPackage ../default.nix { };
  };

in [
  (import "${(builtins.fetchGit {
    url = "https://github.com/tweag/gomod2nix.git";
    rev = "f8ad3b8024896b3c7f571f068c168643708822de";
  })}/overlay.nix")
  (import "${(builtins.fetchGit {
    url = "https://github.com/nix-community/poetry2nix.git";
    rev = "31e3d16d65e65a11ae8e1aaeb9c2e7748144617f";
  })}/overlay.nix")
  overlay
]
