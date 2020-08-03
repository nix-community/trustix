let
  pkgs = import ./nix;

in pkgs.mkShell {

  buildInputs = [
    pkgs.hivemind  # Process monitoring in development
    pkgs.reflex  # File watcher utility
    pkgs.go
    pkgs.vgo2nix

    # Data store
    pkgs.mariadb
    pkgs.trillian
  ];

  TRILLIAN_SCHEMA = "${pkgs.trillian.src}/storage/mysql/schema/storage.sql";

  MYSQL_ROOT_PASSWORD = "trustix-test-log";
  MYSQL_DATABASE = "test";
  MYSQL_USER = "test";
  MYSQL_PASSWORD = "trustix-test-log";

  shellHook = ''
    unset GOPATH
  '';

}
