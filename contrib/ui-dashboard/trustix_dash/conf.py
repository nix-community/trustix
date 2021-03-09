from pydantic import BaseModel
import os


class SettingsModel(BaseModel):
    trustix_rpc: str = os.environ["TRUSTIX_RPC"]
    default_attr: str = "hello.x86_64-linux"
    binary_cache_proxy: str = os.environ["TRUSTIX_BINARY_CACHE_PROXY"]


settings = SettingsModel()
