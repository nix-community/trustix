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

from __future__ import annotations
from trustix_nix_reprod import models as db_models
from pynixutil import b32encode
from typing import (
    Dict,
    List,
    Type,
    Set,
    Any,
)
from pydantic import BaseModel as _BaseModel
import codecs
import json


__all__ = (
    "DerivationReproducibilityStats",
    "DerivationReproducibility",
)


def _json_dumps(v, *, default) -> str:
    return json.dumps(v, default=default, option=json.OPT_NON_STR_KEYS).decode()


class BaseModel(_BaseModel):
    class Config:
        json_loads = json.loads
        json_dumps = _json_dumps


class DerivationOutputResult(BaseModel):
    output_id: str
    output_hash: str
    log_id: int

    @classmethod
    def from_db(cls, mdl: db_models.DerivationOutputResult) -> DerivationOutputResult:
        return cls(
            output_id=b32encode(mdl.output_id),  # type: ignore
            output_hash=codecs.encode(mdl.output_hash, "hex").decode(),
            log_id=mdl.log_id,  # type: ignore
        )


PATH_T = Dict[str, Dict[str, Dict[str, List[DerivationOutputResult]]]]


class Log(BaseModel):
    id: int
    name: str
    tree_size: int

    @classmethod
    def from_db(cls, mdl: db_models.Log) -> Log:
        return cls(
            id=mdl.id,  # type: ignore
            name=mdl.name,
            tree_size=mdl.tree_size,
        )


class DerivationReproducibilityStats(BaseModel):
    pct_reproduced: float
    num_reproduced: int
    num_outputs: int


class DerivationReproducibility(BaseModel):

    unreproduced_paths: PATH_T
    reproduced_paths: PATH_T
    unknown_paths: PATH_T
    missing_paths: PATH_T
    drv_path: str

    logs: Dict[int, Log]

    statistics: DerivationReproducibilityStats

    class Config:
        @staticmethod
        def schema_extra(
            schema: Dict[str, Any], model: Type["DerivationReproducibility"]
        ) -> None:

            massage_keys: Set[str] = set(
                [
                    "unreproduced_paths",
                    "reproduced_paths",
                    "unknown_paths",
                    "missing_paths",
                ]
            )

            key_descr: List[str] = ["drv", "output", "output_hash"]

            for key, prop in schema.get("properties", {}).items():
                if key not in massage_keys:
                    continue

                p = prop
                for title in key_descr:
                    p = p["additionalProperties"]
                    p["title"] = title


class SearchResponse(BaseModel):
    derivations_by_attr: Dict[str, Set[str]]


class SuggestResponse(BaseModel):
    attrs: List[str]


class AttrsResponse(BaseModel):
    attr_stats: Dict[str, Dict[str, DerivationReproducibility]]


class DerivationOutputResultsUniqueResponse(BaseModel):
    results: List[DerivationOutputResult]


class DiffResponse(BaseModel):
    narinfo: Dict[str, Dict]
    diffoscope: Dict
