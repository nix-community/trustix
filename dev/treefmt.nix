{
  # Used to find the project root
  projectRootFile = "flake.lock";

  programs.nixpkgs-fmt.enable = true;
  programs.gofmt.enable = true;

  programs.clang-format.enable = true;
  settings.formatter.clang-format.includes = [
    "*.proto"
  ];
}
