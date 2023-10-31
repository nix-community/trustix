{ pkgs, system, lib, packages }:

let
  testing = import "${pkgs.path}/nixos/lib/testing-python.nix" { inherit pkgs system; };

in
testing.makeTest {
  name = "trustix";
  meta = {
    maintainers = [ lib.maintainers.adisbladis ];
  };

  nodes.machine = { ... }: {
    imports = [ ./default.nix ];

    # Hack around pkgs not working as intended
    nixpkgs.overlays = [
      (_: _: packages)
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
