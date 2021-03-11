from __future__ import annotations
from trustix_dash import models as db_models
from pynix import b32encode
from typing import (
    Dict,
    List,
)
from pydantic import BaseModel
import codecs


__all__ = (
    "DerivationReproducibilityStats",
    "DerivationReproducibility",
)


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

    # # TODO Log -> Pydantic
    logs: Dict[int, Log]

    statistics: DerivationReproducibilityStats
