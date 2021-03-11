let

  overlay = self: super: {
    trustix = self.callPackage ../default.nix { };
  };

in [
  (import "${(builtins.fetchGit {
    url = "https://github.com/tweag/gomod2nix.git";
    rev = "f8ad3b8024896b3c7f571f068c168643708822de";
  })}/overlay.nix")
  (self: super: {
    npmlock2nix = import "${(builtins.fetchGit {
      url = "https://github.com/tweag/npmlock2nix.git";
      rev = "7a321e2477d1f97167847086400a7a4d75b8faf8";
    })}/default.nix" { pkgs = self; };
  })
  overlay
]
