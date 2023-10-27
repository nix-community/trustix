{ npmlock2nix, lib, gitignoreSource }:

npmlock2nix.v1.build {
  src = gitignoreSource ./.;
  installPhase = "cp -r dist $out";
  buildCommands = [ "npm run build" ];
}
