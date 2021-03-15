from trustix_nix_reprod.models import (
    DerivationAttr,
)
from trustix_nix_reprod.lib import (
    flatten,
)
from typing import (
    List,
    Dict,
    Set,
)
import Levenshtein  # type: ignore


__all__ = (
    "search_derivations",
    "suggest_attrs",
)


async def search_derivations(term: str) -> Dict[str, Set[str]]:
    if len(term) < 3:
        raise ValueError("Search term too short")

    results = await DerivationAttr.filter(attr__startswith=term)

    derivations_by_attr: Dict[str, Set[str]] = {}
    for result in results:
        derivation_id: str = result.derivation_id  # type: ignore
        derivations_by_attr.setdefault(result.attr, set()).add(derivation_id)

    return derivations_by_attr


async def suggest_attrs(attr_prefix: str) -> List[str]:
    if len(attr_prefix) < 3:
        raise ValueError("Prefix too short")

    return sorted(
        flatten(
            await DerivationAttr.filter(attr__startswith=attr_prefix)
            .distinct()
            .values_list("attr")
        ),
        key=lambda x: Levenshtein.ratio(attr_prefix, x),
        reverse=True,
    )
