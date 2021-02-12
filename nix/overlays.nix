let

  overlay = self: super: {
    trustix = self.callPackage ../default.nix { };
  };

in [
  (import "${(builtins.fetchGit {
    url = "https://github.com/tweag/gomod2nix.git";
    rev = "f8ad3b8024896b3c7f571f068c168643708822de";
  })}/overlay.nix")
  overlay
]
