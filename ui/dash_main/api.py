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
import grpc
from django.conf import settings
import typing
import pynix


channel = grpc.insecure_channel(settings.TRUSTIX_RPC)
stub = trustix_pb2_grpc.TrustixCombinedRPCStub(channel)


def check_reprod(attr: str):
    # outputs = DerivationOutputResult.objects.filter(output__derivation__attr=attr)

    print(Derivation.objects.filter(attr="hello").values_list("refs", flat=True))

    # drv = Derivation.objects.get(attr=attr)
    # print(drv.refs.all())

    # def get_refs(drv):
    #     for ref in drv.refs.all():
    #         print(ref)
    #         get_refs(ref)

    # get_refs(drv)


@transaction.atomic
def index_eval(commit_sha: str):

    # gh = Github()
    # repo = gh.get_repo("NixOS/nixpkgs")
    # commit = repo.get_commit(commit_sha)

    try:
        evaluation = Evaluation.objects.get(commit=commit_sha)
    except Evaluation.DoesNotExist:
        evaluation = Evaluation.objects.create(
            commit=commit_sha,
            # timestamp=dateutil.parser.parse(
            #     commit.raw_data["commit"]["committer"]["date"]
            # ),
            timestamp=datetime.utcnow(),
        )

    seen: typing.Set[str] = set()

    def gen_drvs(
        attr: typing.Optional[str], drv_path: str
    ) -> typing.Generator[typing.Tuple[typing.Optional[str], str, typing.Any], None, None]:
        with open(drv_path) as f:
            drv = pynix.drvparse(f.read())

        yield (attr, drv_path, drv)
        seen.add(drv_path)

        for drvpath in drv["inputDrvs"]:
            if drvpath not in seen:
                yield from gen_drvs(None, drvpath)

    def gen_drvs_attrs() -> typing.Generator[typing.Tuple[typing.Optional[str], str, typing.Any], None, None]:
        for attr, pkg in ijson.kvitems(open("./output"), ""):
            if "error" in pkg:
                continue

            print(attr)
            attr = attr.rsplit(".", 1)[0]
            yield from gen_drvs(attr, pkg["drvPath"])

    for attr, drv_path, drv in gen_drvs_attrs():
        print(drv_path)

        d, created = Derivation.objects.get_or_create(
            drv=drv_path.split("/")[-1]
        )
        if created:
            d.system = drv["platform"]
        if not d.attr and attr:
            d.attr = attr

        for ref in drv["inputDrvs"]:
            d.refs.add(ref.split("/")[-1])

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
