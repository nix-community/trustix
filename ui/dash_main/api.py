from dash_main.models import (
    Log,
    Evaluation,
    Derivation,
    DerivationOutput,
    DerivationOutputResult,
)
from django.core.paginator import Paginator
from django.db import transaction
from github import Github
import dateutil.parser
import ijson  # type: ignore
from proto import trustix_pb2_grpc
from proto import trustix_pb2
import grpc
import typing
from django.conf import settings

from .util import decode_nix_b32, get_drv_refs


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
def index_eval_outputs(commit_sha: str):
    channel = grpc.insecure_channel(settings.TRUSTIX_RPC)
    stub = trustix_pb2_grpc.TrustixCombinedRPCStub(channel)

    batch_size = 500

    def gen_requests():
        qs = DerivationOutput.objects.all()
        qs = qs.filter(derivation__evaluations__commit=commit_sha)
        qs = qs.order_by("input_hash")
        for p in Paginator(qs, batch_size):
            for drv_output in p.object_list:
                yield trustix_pb2.KeyRequest(Key=drv_output.input_hash)

    existing_logs: typing.Set[str] = set()

    for idx, resp in enumerate(stub.GetStream(gen_requests())):
        print(f"Indexing: {idx}")

        if len(resp.Entries) == 0:
            continue

        for logname, entry in resp.Entries.items():
            if logname not in existing_logs:
                Log.objects.get_or_create(name=logname)
                existing_logs.add(logname)

            output, _ = DerivationOutput.objects.get_or_create(input_hash=resp.Key)

            try:
                DerivationOutputResult.objects.get(output_id=output.id, log_id=logname)
            except DerivationOutputResult.DoesNotExist:
                DerivationOutputResult.objects.create(
                    output_id=output.id, log_id=logname, output_hash=entry.Value
                )
