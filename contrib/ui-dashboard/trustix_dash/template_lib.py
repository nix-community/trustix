import urllib.parse


def quote_drv_url(drv_path: str) -> str:
    return urllib.parse.quote(urllib.parse.quote(drv_path, safe=""))
