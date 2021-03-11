from trustix_dash.api.models import DerivationOutputResult
from trustix_dash import models as db_models
from typing import (
    Dict,
    List,
)


__all__ = ("get_derivation_output_results_unique",)


async def get_derivation_output_results_unique(
    *output_hash: bytes,
) -> List[DerivationOutputResult]:
    if not output_hash:
        return []

    results: Dict[bytes, db_models.DerivationOutputResult] = {
        result.output_hash: result  # type: ignore
        for result in await db_models.DerivationOutputResult.filter(
            output_hash__in=output_hash
        )
    }

    if len(results) != len(output_hash):
        raise ValueError(
            "{} ids passed but only returned {} results".format(
                len(output_hash), len(results)
            )
        )

    return [DerivationOutputResult.from_db(results[h]) for h in output_hash]
