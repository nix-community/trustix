from __future__ import annotations
from trustix_nix_reprod import models as db_models
from pynix import b32encode
from typing import (
    Dict,
    List,
    Type,
    Set,
    Any,
)
from pydantic import BaseModel as _BaseModel
import codecs
import orjson


__all__ = (
    "DerivationReproducibilityStats",
    "DerivationReproducibility",
)


def _orjson_dumps(v, *, default) -> str:
    return orjson.dumps(v, default=default).decode()


class BaseModel(_BaseModel):

    class Config:
        json_loads = orjson.loads
        json_dumps = _orjson_dumps


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
