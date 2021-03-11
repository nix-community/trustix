from typing import (
    Coroutine,
    Dict,
    List,
    Set,
)
from trustix_dash.models import (
    DerivationOutputResult,
    DerivationOutput,
    Derivation,
    Log,
)
from tortoise.query_utils import Q
import asyncio
import codecs


__all__ = ("get_derivation_reproducibility",)


async def _get_derivation_outputs(drv: str) -> List[Derivation]:
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

    items: List[Derivation] = []
    for items_ in await asyncio.gather(*coros):
        items.extend(items_)

    return items


async def get_derivation_reproducibility(drv_path: str):

    drvs = await _get_derivation_outputs(drv_path)

    # Description: drv -> output -> output_hash -> List[result]
    PATH_T = Dict[str, Dict[str, Dict[str, List[DerivationOutputResult]]]]
    unreproduced_paths: PATH_T = {}
    reproduced_paths: PATH_T = {}
    # Paths only built by one log
    unknown_paths: PATH_T = {}
    # Paths not built by any known log
    missing_paths: PATH_T = {}

    log_ids: Set[int] = set()

    def append_output(paths_d: PATH_T, drv: Derivation, output: DerivationOutput):
        """Append an output to the correct datastructure (paths_d)"""
        current = paths_d.setdefault(drv.drv, {}).setdefault(output.output, {})
        for result in output.derivationoutputresults:  # type: ignore
            current.setdefault(
                codecs.encode(result.output_hash, "hex").decode(), []
            ).append(result)
            log_ids.add(result.log_id)

    num_outputs: int = 0

    for drv in drvs:
        output: DerivationOutput
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
        log.id: log for log in await Log.filter(id__in=log_ids)  # type: ignore
    }

    num_reproduced = sum(len(v) for v in reproduced_paths.values())

    # TODO: Get first evaluation timestamp

    return {
        "unreproduced_paths": unreproduced_paths,
        "reproduced_paths": reproduced_paths,
        "unknown_paths": unknown_paths,
        "missing_paths": missing_paths,
        "drv_path": drv_path,
        "logs": logs,
        "statistics": {
            "pct_reproduced": round(num_outputs / 100 * num_reproduced, 2),
            "num_reproduced": num_reproduced,
            "num_outputs": num_outputs,
        },
    }
