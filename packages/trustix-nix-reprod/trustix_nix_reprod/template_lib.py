from diffoscope.presenters.html import HTMLPresenter  # type: ignore
from diffoscope.readers.json import JSONReaderV1  # type: ignore
from markupsafe import Markup
from unittest import mock
from typing import Dict
import urllib.parse
import contextlib
import json


__all__ = (
    "drv_url_quote",
    "json_render",
    "js_url",
    "diffoscope_render",
)


def drv_url_quote(drv_path: str) -> str:
    return urllib.parse.quote(urllib.parse.quote(drv_path, safe=""))


def json_render(x) -> str:
    return json.dumps(x)


def js_url(filename: str) -> Markup:
    return Markup("/".join(("", "js", filename)))


@contextlib.contextmanager
def _make_diffoscope_printer(callback):
    def fn(html: str):
        callback(html)
    yield fn


def diffoscope_render(value: Dict) -> Markup:
    """Convert diffoscope JSON output to HTML"""

    s = ""

    def callback(html):
        nonlocal s
        s = html

    diff = JSONReaderV1().load_rec(value)

    with mock.patch("diffoscope.presenters.html.html.make_printer", _make_diffoscope_printer):
        HTMLPresenter().output_html(callback, diff)

    return Markup(s)
