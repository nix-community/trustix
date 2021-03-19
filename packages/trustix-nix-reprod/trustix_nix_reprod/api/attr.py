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


__all__ = (
    "get_attrs_reproducibility",
    "get_attr_reproducibility",
)


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
