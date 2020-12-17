from dash_main.models import (
    Log,
    Evaluation,
    Derivation,
    DerivationOutput,
    DerivationOutputResult,
)
from django.db import transaction
from github import Github
import dateutil.parser
import ijson  # type: ignore
from proto import trustix_pb2_grpc
from proto import trustix_pb2
import grpc  # type: ignore
from django.conf import settings
from functools import lru_cache
from django.db.models import Q
from concurrent import futures
import typing
import pynix


channel = grpc.insecure_channel(settings.TRUSTIX_RPC)
stub = trustix_pb2_grpc.TrustixCombinedRPCStub(channel)


def get_derivation_outputs(drv: typing.Optional[str]) -> typing.List[DerivationOutputResult]:

    items = []
    with futures.ThreadPoolExecutor(max_workers=1) as e:
        f = e.submit(lambda: None)
        print(f.result())

    for q_filter in (Q(refs_all=drv), Q(drv=drv)):
        items.extend(Derivation.objects.filter(q_filter).select_related("output").select_related("output__result"))
    return []


def get_derivation_outputs__(drv: typing.Optional[str]) -> typing.List[DerivationOutputResult]:
    def select_related(qs):
        return qs.select_related("output").select_related("output__derivation")

    items = list(
        select_related(
            DerivationOutputResult.objects.filter(Q(output__derivation__refs_all=drv))
        )
    )
    items.extend(
        list(
            DerivationOutputResult.objects.exclude(id__in=[i.id for i in items]).filter(
                output__derivation__drv=drv
            )
        )
    )

    return items


@transaction.atomic
def index_eval(commit_sha: str):

    gh = Github()
    repo = gh.get_repo("NixOS/nixpkgs")
    commit = repo.get_commit(commit_sha)

    try:
        evaluation = Evaluation.objects.get(commit=commit_sha)
    except Evaluation.DoesNotExist:
        evaluation = Evaluation.objects.create(
            commit=commit_sha,
            timestamp=dateutil.parser.parse(
                commit.raw_data["commit"]["committer"]["date"]
            ),
        )

    refs: typing.Dict[str, typing.Set[str]] = {}

    # Consider:
    # - Fast short circuit when drv is already indexed

    @lru_cache(maxsize=30_000)
    def drv_read(drv_path: str) -> typing.Dict:
        with open(drv_path) as f:
            return pynix.drvparse(f.read())

    def gen_drvs(
        attr: typing.Optional[str], drv_path: str
    ) -> typing.Generator[
        typing.Tuple[
            typing.Optional[str], typing.Dict, typing.Set[str], typing.Set[str], str
        ],
        None,
        None,
    ]:
        drv = drv_read(drv_path)

        # Direct dependencies
        refs_direct: typing.Set[str] = set(drv["inputDrvs"])

        # All dependencies (recursive, flattened)
        refs_all = refs_direct.copy()

        for input_ in drv["inputDrvs"]:
            if input_ not in refs:
                yield from gen_drvs(None, input_)

            # If the input _still_ doesn't exist it means it's a fixed-output
            # and should be filtered out
            try:
                refs_all = refs_all | refs[input_]
            except KeyError:
                refs_direct.remove(input_)
                refs_all.remove(input_)

        # Filter fixed outputs
        if all("hashAlgo" in d for d in drv["outputs"].values()):
            return

        refs[drv_path] = refs_direct

        yield (attr, drv, refs_direct, refs_all, drv_path)

    def gen_drvs_attrs() -> typing.Generator[
        typing.Tuple[
            typing.Optional[str], typing.Dict, typing.Set[str], typing.Set[str], str
        ],
        None,
        None,
    ]:
        for attr, pkg in ijson.kvitems(open("./output"), ""):
            if "error" in pkg:
                continue

            attr = attr.rsplit(".", 1)[0]
            yield from gen_drvs(attr, pkg["drvPath"])

    for (attr, drv, refs_direct, refs_all, drv_path) in gen_drvs_attrs():
        if attr:
            print(f"Indexing {attr}")

        d, created = Derivation.objects.get_or_create(drv=drv_path.split("/")[-1])
        changed = False
        if created:
            changed = True
            d.system = drv["platform"]

            for ref in refs_direct:
                d.refs_direct.add(ref.split("/")[-1])

            for ref in refs_all:
                d.refs_all.add(ref.split("/")[-1])

        if not d.attr and attr:
            changed = True
            d.attr = attr

        if changed:
            d.save()

        for output, store_path in drv["outputs"].items():
            store_path = store_path["path"].split("/")[-1]
            input_hash = pynix.b32decode(store_path.split("-", 1)[0])

            try:
                DerivationOutput.objects.get(input_hash=input_hash)
            except DerivationOutput.DoesNotExist:
                DerivationOutput.objects.create(
                    derivation=d,
                    input_hash=input_hash,
                    output=output,
                    store_path=store_path,
                )


@transaction.atomic
def index_log(log, sth):
    start = log.tree_size
    finish = sth.TreeSize - 1

    chunks = list(range(log.tree_size, finish, 500))
    if chunks[-1] != finish:
        chunks.append(finish)

    start = chunks[0]
    for finish in chunks:
        resp = stub.GetLogEntries(
            trustix_pb2.GetLogEntriesRequestNamed(
                LogName=log.name,
                Start=start,
                Finish=finish,
            )
        )
        print(f"Indexing logname={log.name}, start={start}, finish={finish}")

        for leaf in resp.Leaves:
            try:
                DerivationOutputResult.objects.get(output_id=leaf.Key, log=log)
            except DerivationOutputResult.DoesNotExist:
                DerivationOutputResult.objects.create(
                    output_id=leaf.Key, output_hash=leaf.Value, log=log
                )

        start = finish


def index_logs():
    req = trustix_pb2.LogsRequest()
    resp = stub.Logs(req)
    for logname, meta in resp.Logs.items():
        try:
            log = Log.objects.get(name=logname)
        except Log.DoesNotExist:
            log = Log.objects.create(name=logname, tree_size=0)

        index_log(log, meta.STH)
