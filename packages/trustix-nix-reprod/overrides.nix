{ pkgs }: self: super: {

  # This package is a rust thing with Cargo.lock/Cargo.toml
  orjson = super.orjson.override {
    preferWheel = true;
  };

}
