{ npmlock2nix, lib, gitignoreSource }:

npmlock2nix.build {
  src = gitignoreSource ./.;
  installPhase = "cp -r dist $out";
  buildCommands = [ "npm run build" ];
}
