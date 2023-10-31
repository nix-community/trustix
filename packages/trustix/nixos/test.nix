{ pkgs, lib }:

let
  testing = import "${pkgs.path}/nixos/lib/testing-python.nix" { inherit pkgs system; };

in
testing.makeTest {

  inherit pkgs;

  name = "trustix";
  meta = {
    maintainers = [ lib.maintainers.adisbladis ];
  };

  machine = { ... }: {
    imports = [ ./default.nix ];

    # Hack around pkgs not working as intended
    nixpkgs.overlays = [
      (self: super: {
        inherit (pkgs) trustix;
      })
    ];

    services.trustix.enable = true;
  };

  testScript =
    ''
      start_all()
      machine.wait_for_unit("trustix")
      machine.shutdown()
    '';

}
