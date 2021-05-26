{ pkgs ? import ./nix }:

{
  packages = {
    inherit (pkgs) trustix trustix-doc trustix-nix trustix-nix-reprod;
  };

  containers =
    let
      inheritContainer = attr: {}@args: pkgs.dockerTools.buildLayeredImage
        {
          inherit (pkgs.${attr}) name;
          contents = [ pkgs.${attr} ];
        } // args;
    in
    {
      trustix = inheritContainer "trustix" { };
      trustix-nix = inheritContainer "trustix-nix" { };
      trustix-nix-reprod = inheritContainer "trustix-nix-reprod" { };
    };

}
