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

from pydantic import BaseModel
from collections import OrderedDict
from typing import (
    Dict,
    List,
    Any,
)


class DiffEntry(BaseModel):
    value_a: Any
    value_b: Any

    @property
    def has_diff(self) -> bool:
        return self.value_a != self.value_b


def json_diff(a: Dict, b: Dict):

    keys: List[str] = []
    for k in list(a.keys()) + list(b.keys()):
        if k not in keys:
            keys.append(k)

    ret: Dict[str, Any] = OrderedDict()
    for key in keys:
        ret[key] = DiffEntry(value_a=a.get(key), value_b=b.get(key))

    return ret
