{ stdenv
, lib
, buildGoModule
, fetchFromGitHub
}:

let
  pname = "trillian";
  version = "1.3.2";

in buildGoModule {
  inherit pname version;

  modSha256 = "03k3pn5y91gi06ks8v7sj1syyw01py0m8rz1rhyj2fwc3lycdb15";

  src = fetchFromGitHub {
    owner = "google";
    repo = pname;
    rev = "v${version}";
    sha256 = "03zmbyncr28vaa319y9v8ahn0x2397ssa4slqdv0gyb8jwcy59b7";
  };

  subPackages = [
    "server/trillian_log_server"
    "server/trillian_log_signer"
    "server/trillian_map_server"
  ];

  meta = with lib; {
    homepage = "https://github.com/google/trillian";
    description = "A transparent, highly scalable and cryptographically verifiable data store.";
    license = [ licenses.asl20 ];
    maintainers = [ maintainers.adisbladis ];
  };
}
