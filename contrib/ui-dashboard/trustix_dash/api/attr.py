from trustix_dash.api.derivation import (
    get_derivation_reproducibility,
)
from trustix_dash.models import (
    Derivation,
)
from trustix_dash.lib import (
    flatten,
    unique,
)
from collections import OrderedDict
import asyncio


async def get_attr_reproducibility(attr: str):
    drvs = list(
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
