from trustix_dash.models import DerivationOutputResult
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

    results: Dict[bytes, DerivationOutputResult] = {
        result.output_hash: result  # type: ignore
        for result in await DerivationOutputResult.filter(output_hash__in=output_hash)
    }

    if len(results) != len(output_hash):
        raise ValueError(
            "{} ids passed but only returned {} results".format(
                len(output_hash), len(results)
            )
        )

    return [results[h] for h in output_hash]
