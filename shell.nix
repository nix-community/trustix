let
  pkgs = import ./nix;

  pythonEnv = pkgs.python3.withPackages(ps: []);

in pkgs.mkShell {

  buildInputs = [
    pkgs.hivemind  # Process monitoring in development
    pkgs.reflex  # File watcher utility
    pkgs.go
    pkgs.vgo2nix

    pkgs.protobuf

    # For development scripts
    pythonEnv

    # Data store
    pkgs.mariadb
    pkgs.trillian
  ];

  TRILLIAN_SCHEMA = "${pkgs.trillian.src}/storage/mysql/schema/storage.sql";
  MYSQL_BASEDIR = pkgs.mysql;

  MYSQL_ROOT_PASSWORD = "trustix-test-log";
  MYSQL_DATABASE = "test";
  MYSQL_USER = "test";
  MYSQL_PASSWORD = "trustix-test-log";

  shellHook = ''
    unset GOPATH
  '';

}
