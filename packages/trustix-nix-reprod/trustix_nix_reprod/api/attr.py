from trustix_nix_reprod.api.derivation import (
    get_derivation_reproducibility,
)
from trustix_nix_reprod.api.models import (
    DerivationReproducibility,
    AttrsResponse,
)
from trustix_nix_reprod.models import (
    Derivation,
)
from trustix_nix_reprod.lib import (
    flatten,
    unique,
)
from collections import OrderedDict
import asyncio
from typing import (
    List,
    Dict,
)


__all__ = ("get_attr_reproducibility",)


async def get_attr_reproducibility(attr: str) -> Dict[str, DerivationReproducibility]:
    drvs: List[str] = list(
        unique(
            flatten(
                await Derivation.filter(derivationattrs__attr=attr)
                .limit(100)
                .order_by("derivationevals__eval__timestamp")
                .values_list("drv")
            )
        )
    )
    return OrderedDict(
        zip(
            drvs,
            await asyncio.gather(
                *[get_derivation_reproducibility(drv) for drv in drvs]
            ),
        )
    )


async def get_attrs_reproducibility(
    attrs: List[str],
) -> AttrsResponse:
    return AttrsResponse(
        attr_stats=OrderedDict(
            zip(
                attrs,
                await asyncio.gather(*[get_attr_reproducibility(attr) for attr in attrs]),
            )
        ),
    )
