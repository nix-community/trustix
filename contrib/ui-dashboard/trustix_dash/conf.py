from pydantic import BaseModel
import typing
import os


class SettingsModel(BaseModel):
    trustix_rpc: str = os.environ["TRUSTIX_RPC"]
    binary_cache_proxy: str = os.environ["TRUSTIX_BINARY_CACHE_PROXY"]
    db_uri: str = os.environ["DB_URI"]

    default_attrs: typing.List[str] = (
        os.environ["DEFAULT_ATTRS"].split(":")
        if "DEFAULT_ATTRS" in os.environ
        else [
            "hello.x86_64-linux",
        ]
    )

    placeholder_attr: str = "hello.x86_64-linux"

    supported_systems: typing.List[str] = (
        os.environ["SUPPORTED_SYSTEMS"].split(":")
        if "SUPPORTED_SYSTEMS" in os.environ
        else [
            "aarch64-linux",
            "x86_64-linux",
            "x86_64-darwin",
        ]
    )

    @property
    def tortoise_config(self) -> typing.Dict:
        return {
            "connections": {
                "default": self.db_uri,
            },
            "apps": {
                "trustix_dash": {
                    "models": ["trustix_dash.models", "aerich.models"],
                }
            },
            "use_tz": False,
            "timezone": "UTC",
        }


settings = SettingsModel()


TORTOISE_ORM = settings.tortoise_config
