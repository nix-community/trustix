from trustix_nix_reprod.api.models import (
    DerivationOutputResultsUniqueResponse,
    DerivationOutputResult,
)
from trustix_nix_reprod.cache import cached
from trustix_nix_reprod.conf import settings
from trustix_nix_reprod import models as db_models
from typing import Dict


__all__ = ("get_derivation_output_results_unique",)


@cached(model=DerivationOutputResultsUniqueResponse, ttl=settings.cache_ttl.diff)
async def get_derivation_output_results_unique(
    *output_hash: bytes,
) -> DerivationOutputResultsUniqueResponse:
    if not output_hash:
        raise ValueError("No hashes provided")

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

    return DerivationOutputResultsUniqueResponse(
        results=[DerivationOutputResult.from_db(results[h]) for h in output_hash]
    )
