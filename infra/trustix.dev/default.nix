{ config, pkgs, ... }:

let
  servedDomains = [
    "trustix.dev"
    "r13y.trustix.dev"
  ];

in
{
  imports = [
    ./hardware-configuration.nix
    ./trustix.dev.nix
  ];

  boot.loader.grub.enable = true;
  boot.loader.grub.version = 2;
  boot.loader.grub.device = "/dev/sda"; # or "nodev" for efi only

  networking.hostName = "trustixdotdev";

  networking.usePredictableInterfaceNames = false;
  networking.dhcpcd.enable = false;
  systemd.network = {
    enable = true;
    networks."ethernet".extraConfig = ''
      [Match]
      Type = ether
      [Network]
      DHCP = ipv4
      Address = 2a01:4f9:c012:7359::1/64
      Gateway = fe80::1
    '';
  };

  services.nginx = {
    enable = true;
    # Set sane TLS defaults for all vhosts
    virtualHosts = builtins.listToAttrs (map
      (n: {
        name = n;
        value = {
          enableACME = true;
          forceSSL = true;
        };
      })
      servedDomains);
  };


  services.openssh = {
    enable = true;
    startWhenNeeded = true;
  };

  users.extraUsers.root.openssh.authorizedKeys.keys = [
    "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCtr+rcxCZBAAqt8ocvhEEdBWfnRBCljjQPtC6Np24Y3H/HMe3rugsu3OhPscRV1k5hT+UlA2bpN8clMFAfK085orYY7DMUrgKQzFB7GDnOvuS1CqE1PRw7/OHLcWxDwf3YLpa8+ZIwMHFxR2gxsldCLGZV/VukNwhEvWs50SbXwVrjNkwA9LHy3Or0i6sAzU711V3B2heB83BnbT8lr3CKytF3uyoTEJvDE7XMmRdbvZK+c48bj6wDaqSmBEDrdNncsqnReDjScdNzXgP1849kMfIUwzXdhEF8QRVfU8n2A2kB0WRXiGgiL4ba5M+N9v1zLdzSHcmB0veWGgRyX8tN cardno:FF7F00" # adisbladis
  ];

  networking.firewall.allowedTCPPorts = [
    80
    443
  ];
  networking.firewall.allowedUDPPorts = [
    80
    443
  ];

  system.stateVersion = "22.11";
}
