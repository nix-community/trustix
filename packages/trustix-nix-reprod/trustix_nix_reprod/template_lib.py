# Trustix
# Copyright (C) 2021 Tweag IO

# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.

# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.

# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

from diffoscope.presenters.html import HTMLPresenter  # type: ignore
from diffoscope.readers.json import JSONReaderV1  # type: ignore
from markupsafe import Markup
from unittest import mock
from typing import Dict
import urllib.parse
import contextlib
import json
import re


__all__ = (
    "drv_url_quote",
    "json_render",
    "js_url",
    "diffoscope_render",
)


def drv_url_quote(drv_path: str) -> str:
    return urllib.parse.quote(urllib.parse.quote(drv_path, safe=""))


def json_render(x) -> Markup:
    if isinstance(x, list):
        return Markup(", ".join([json.dumps(i, indent=2) for i in x]))

    return Markup(json.dumps(x, indent=2))


def js_url(filename: str) -> Markup:
    return Markup("/".join(("", "js", filename)))


@contextlib.contextmanager
def _make_diffoscope_printer(callback):
    def fn(html: str):
        callback(html)
    yield fn


def diffoscope_render(value: Dict) -> str:
    """Convert diffoscope JSON output to HTML"""

    s = ""

    def callback(html):
        nonlocal s
        s = html

    diff = JSONReaderV1().load_rec(value)

    with mock.patch("diffoscope.presenters.html.html.make_printer", _make_diffoscope_printer):
        HTMLPresenter().output_html(callback, diff)

    # Diffoscope leaks argv into title
    s = re.sub(r"<title>.*?</title>", "<title>Trustix Diffoscope</title>", s, flags=re.DOTALL)

    return s
