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

from .util import decode_nix_b32, get_drv_refs


channel = grpc.insecure_channel(settings.TRUSTIX_RPC)
stub = trustix_pb2_grpc.TrustixCombinedRPCStub(channel)


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

    for attr, pkg in ijson.kvitems(open("./output"), ""):
        if "error" in pkg:
            continue

        attr = attr.rstrip("." + pkg["system"])
        print(attr)

        drv, created = Derivation.objects.get_or_create(
            drv=pkg["drvPath"].split("/")[-1]
        )
        if created:
            drv.system = pkg["system"]
            drv.attr = attr
            for ref in get_drv_refs(pkg["drvPath"]):
                drv.refs.add(ref.split("/")[-1])

        drv.evaluations.add(evaluation)
        drv.save()

        for output, store_path in pkg["outputs"].items():
            store_path = store_path.split("/")[-1]
            DerivationOutput.objects.get_or_create(
                derivation=drv,
                input_hash=decode_nix_b32(store_path.split("-", 1)[0]),
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
        resp = stub.GetLogEntries(trustix_pb2.GetLogEntriesRequestNamed(
            LogName=log.name,
            Start=start,
            Finish=finish,
        ))
        print(f"Indexing logname={log.name}, start={start}, finish={finish}")

        for leaf in resp.Leaves:
            try:
                DerivationOutputResult.objects.get(output_id=leaf.Digest, log=log)
            except DerivationOutputResult.DoesNotExist:
                DerivationOutputResult.objects.create(output_id=leaf.Digest, output_hash=leaf.Value, log=log)

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
