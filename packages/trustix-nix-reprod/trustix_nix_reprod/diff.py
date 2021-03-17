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
