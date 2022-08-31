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

from trustix_nix_reprod.conf import settings
from trustix_nix_reprod.cache import cached
from trustix_nix_reprod.models import (
    DerivationAttr,
)
from trustix_nix_reprod.api.models import (
    SuggestResponse,
    SearchResponse,
)
from trustix_nix_reprod.lib import (
    flatten,
)
from typing import (
    Dict,
    Set,
)
import Levenshtein  # type: ignore


__all__ = (
    "search_derivations",
    "suggest_attrs",
)


@cached(model=SearchResponse, ttl=settings.cache_ttl.search)
async def search_derivations(term: str) -> SearchResponse:

    if len(term) < 3:
        raise ValueError("Search term too short")

    results = await DerivationAttr.filter(attr__startswith=term)

    derivations_by_attr: Dict[str, Set[str]] = {}
    for result in results:
        derivation_id: str = result.derivation_id  # type: ignore
        derivations_by_attr.setdefault(result.attr, set()).add(derivation_id)

    return SearchResponse(derivations_by_attr=derivations_by_attr)


@cached(model=SuggestResponse, ttl=settings.cache_ttl.suggest)
async def suggest_attrs(attr_prefix: str) -> SuggestResponse:
    if len(attr_prefix) < 3:
        raise ValueError("Prefix too short")

    return SuggestResponse(
        attrs=sorted(
            flatten(
                await DerivationAttr.filter(attr__startswith=attr_prefix)
                .distinct()
                .values_list("attr")
            ),
            key=lambda x: Levenshtein.ratio(attr_prefix, x),
            reverse=True,
        ),
    )
