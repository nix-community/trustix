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

from typing import (
    Coroutine,
    Dict,
    List,
    Set,
)
from trustix_nix_reprod.api.models import (
    DerivationReproducibilityStats,
    DerivationReproducibility,
    DerivationOutputResult,
    PATH_T,
    Log,
)
from trustix_nix_reprod import models as db_models
from trustix_nix_reprod.cache import cached
from trustix_nix_reprod.conf import settings
from tortoise.query_utils import Q
import asyncio


__all__ = ("get_derivation_reproducibility",)


async def _get_derivation_outputs(drv: str) -> List[db_models.Derivation]:
    async def filter(q_filter):
        qs = (
            db_models.Derivation.filter(q_filter)
            .prefetch_related("derivationoutputs")
            .prefetch_related("derivationoutputs__derivationoutputresults")
        )
        return await qs

    coros: List[Coroutine] = [
        filter(q_filter)
        for q_filter in (Q(from_ref_recursive__referrer=drv), Q(drv=drv))
    ]

    items: List[db_models.Derivation] = []
    for items_ in await asyncio.gather(*coros):
        items.extend(items_)

    return items


@cached(model=DerivationReproducibility, ttl=settings.cache_ttl.drv_reprod)
async def get_derivation_reproducibility(drv_path: str) -> DerivationReproducibility:

    drvs = await _get_derivation_outputs(drv_path)

    # Description: drv -> output -> output_hash -> List[result]
    unreproduced_paths: PATH_T = {}
    reproduced_paths: PATH_T = {}
    # Paths only built by one log
    unknown_paths: PATH_T = {}
    # Paths not built by any known log
    missing_paths: PATH_T = {}

    log_ids: Set[int] = set()

    def append_output(
        paths_d: PATH_T, drv: db_models.Derivation, output: db_models.DerivationOutput
    ):
        """Append an output to the correct datastructure (paths_d)"""
        current = paths_d.setdefault(drv.drv, {}).setdefault(output.output, {})
        for result in output.derivationoutputresults:  # type: ignore
            res = DerivationOutputResult.from_db(result)
            current.setdefault(res.output_hash, []).append(res)
            log_ids.add(result.log_id)

    num_outputs: int = 0

    for drv in drvs:
        output: db_models.DerivationOutput
        for output in drv.derivationoutputs:  # type: ignore
            output_hashes: Set[bytes] = set(
                result.output_hash for result in output.derivationoutputresults  # type: ignore
            )

            num_outputs += 1

            if not output_hashes:
                append_output(missing_paths, drv, output)

            elif len(output_hashes) == 1 and len(output.derivationoutputresults) > 1:  # type: ignore
                append_output(reproduced_paths, drv, output)

            elif len(output_hashes) == 1:
                append_output(unknown_paths, drv, output)

            elif len(output_hashes) > 1:
                append_output(unreproduced_paths, drv, output)

            else:
                raise RuntimeError("Logic error")

    logs: Dict[int, Log] = {
        log.id: Log.from_db(log) for log in await db_models.Log.filter(id__in=log_ids)  # type: ignore
    }

    num_reproduced = sum(len(v) for v in reproduced_paths.values())

    return DerivationReproducibility(
        unreproduced_paths=unreproduced_paths,
        reproduced_paths=reproduced_paths,
        unknown_paths=unknown_paths,
        missing_paths=missing_paths,
        drv_path=drv_path,
        logs=logs,
        statistics=DerivationReproducibilityStats(
            pct_reproduced=round(num_outputs / 100 * num_reproduced, 2),
            num_reproduced=num_reproduced,
            num_outputs=num_outputs,
        ),
    )
