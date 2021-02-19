import asyncio
from trustix_dash.models import (
    DerivationOutputResult,
    Derivation,
)
from tortoise.query_utils import Q
from typing import (
    Coroutine,
    List,
)


async def evaluation_list() -> List[str]:
    """
    Get a list of default evaluations to show in the UI
    """

    return [
        "eval1",
        "eval2",
    ]


async def channels_list() -> List[str]:
    """
    Get a list of default channels to show in the UI
    """

    return [
        "chan1",
        "chan2",
    ]


# TODO: Verify return type
async def get_derivation_outputs(drv: str) -> List[DerivationOutputResult]:
    async def filter(q_filter):
        qs = (
            Derivation.filter(q_filter)
            .prefetch_related("derivationoutputs")
            .prefetch_related("derivationoutputs__derivationoutputresults")
        )
        return await qs

    coros: List[Coroutine] = [
        filter(q_filter)
        for q_filter in (Q(from_ref_recursive__referrer=drv), Q(drv=drv))
    ]

    items: List[DerivationOutputResult] = []
    for items_ in await asyncio.gather(*coros):
        items.extend(items_)

    return items


async def get_derivation_output_results(*ids: int) -> List[DerivationOutputResult]:
    if not ids:
        return []

    results: List[DerivationOutputResult] = await DerivationOutputResult.filter(
        id__in=ids
    )
    if len(results) != len(ids):
        raise ValueError(
            "{} ids passed but only returned {} results".format(len(ids), len(results))
        )

    return results
