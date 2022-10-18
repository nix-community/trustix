{ buildGoApplication
, lib
, gitignoreSource
, makeWrapper
, nix-eval-jobs
, nix
, diffoscope
}:

let
  runtimeDeps = [
    nix
    nix-eval-jobs
    diffoscope
  ];

in
buildGoApplication {
  pname = "trustix-nix-r13y";
  version = "dev";
  pwd = ./.;
  modules = ./gomod2nix.toml;
  src = lib.cleanSourceWith {
    filter = name: type: ! lib.hasSuffix "tests" name;
    src = gitignoreSource ./.;
  };

  nativeBuildInputs = [
    makeWrapper
  ];

  postInstall = ''
    wrapProgram "$out/bin/trustix-nix-r13y" --prefix PATH : "${lib.makeBinPath runtimeDeps}"
  '';
}
