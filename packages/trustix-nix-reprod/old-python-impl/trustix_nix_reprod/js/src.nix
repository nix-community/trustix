{ lib }:

lib.cleanSourceWith {
  filter = name: type: !(lib.hasSuffix "node_modules" name) && !(lib.hasSuffix ".direnv" name) && !(lib.hasSuffix "dist" name);
  src = lib.cleanSource ./.;
}
