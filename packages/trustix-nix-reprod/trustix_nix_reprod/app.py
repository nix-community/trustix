from trustix_python.api import api_pb2
from fastapi.templating import (
    Jinja2Templates,
)
from fastapi.responses import (
    RedirectResponse,
    HTMLResponse,
)
from fastapi import (
    FastAPI,
    Request,
    Form,
)
from typing import (
    Callable,
    TypeVar,
    Union,
    Dict,
    List,
)
from trustix_nix_reprod import (
    template_lib,
    on_startup,
    on_shutdown,
)
import urllib.parse
import requests
import tempfile
import asyncio
import os.path
import codecs
import shlex
import json

from trustix_nix_reprod.staticfiles import StaticFiles

from pydantic import BaseModel as BaseModel

from trustix_nix_reprod.api import (
    get_derivation_output_results_unique,
    get_derivation_reproducibility,
    get_attrs_reproducibility,
    search_derivations,
    suggest_attrs,
)
from trustix_nix_reprod.api.models import (
    DerivationReproducibility,
    SuggestResponse,
    SearchResponse,
    AttrsResponse,
)

from trustix_nix_reprod.proto import (
    get_combined_rpc,
)

from trustix_nix_reprod.conf import settings


SCRIPT_DIR = os.path.dirname(__file__)


app = FastAPI()
app.mount(
    "/static", StaticFiles(directory=os.path.join(SCRIPT_DIR, "static")), name="static"
)
app.mount("/js", StaticFiles(directory=settings.js_store), name="js")


templates = Jinja2Templates(directory=os.path.join(SCRIPT_DIR, "templates"))
templates.env.globals["drv_url_quote"] = template_lib.drv_url_quote
templates.env.globals["json_render"] = template_lib.json_render
templates.env.globals["url_reverse"] = app.url_path_for
templates.env.globals["js_url"] = template_lib.js_url


@app.on_event("startup")
async def startup_event():
    await on_startup()


@app.on_event("shutdown")
async def shutdown_event():
    await on_shutdown()


def make_context(
    request: Request,
    model: BaseModel,
    title: str = "",
) -> Dict:

    ctx = {
        "request": request,
        "title": "Trustix R13Y" + (" ".join((" - ", title)) if title else ""),
        "drv_placeholder": settings.placeholder_attr,
        "model": model,
    }

    return ctx


T = TypeVar("T")


def render_model(
    request: Request, template: str, ctx_callback: Callable[[], Dict], model: T
) -> Union[HTMLResponse, T]:
    """
    Polymorphic response depending on the accept header
    In the case of text/html we render the template and otherwise we use a JSON response
    """
    accept_html: bool = "text/html" in [
        mime.split(";")[0] for mime in request.headers["accept"].split(",")
    ]
    if accept_html:
        return templates.TemplateResponse(template, ctx_callback())
    else:
        return model


@app.get("/", response_class=HTMLResponse, include_in_schema=False)
async def index(request: Request):
    return await attrs(request, ",".join(settings.default_attrs))


@app.get("/drv/{drv_path}", response_model=DerivationReproducibility)
async def drv(request: Request, drv_path: str):
    data = await get_derivation_reproducibility(urllib.parse.unquote(drv_path))
    return render_model(
        request,
        "drv.jinja2",
        lambda: make_context(request, data),
        data,
    )


@app.get("/attrs/{attrs}", response_model=AttrsResponse)
async def attrs(request: Request, attrs: str):
    model = await get_attrs_reproducibility(settings.default_attrs)
    return render_model(
        request,
        "attrs.jinja2",
        lambda: make_context(request, model),
        model,
    )


@app.post("/search_form/", include_in_schema=False, response_class=RedirectResponse)
async def search_form(request: Request, term: str = Form(...)):
    return RedirectResponse(app.url_path_for("search", term=term), status_code=303)


@app.get("/search/{term}", response_model=SearchResponse)
async def search(request: Request, term: str):
    data = await search_derivations(term)
    return render_model(
        request,
        "search.jinja2",
        lambda: make_context(request, data),
        data,
    )


@app.get("/suggest/{attr_prefix}", response_model=SuggestResponse)
async def suggest(request: Request, attr_prefix: str):
    return await suggest_attrs(attr_prefix)


@app.post("/diff_form/", response_class=RedirectResponse, include_in_schema=False)
async def diff_form(request: Request, output_hash: List[str] = Form(...)):

    if len(output_hash) < 1:
        raise ValueError("Need at least 2 entries to diff")
    if len(output_hash) > 2:
        raise ValueError("Received more than 2 entries to diff")

    return RedirectResponse(
        app.url_path_for(
            "diff", output_hash_1_hex=output_hash[0], output_hash_2_hex=output_hash[1]
        ),
        status_code=303,
    )


@app.get("/diff/{output_hash_1_hex}/{output_hash_2_hex}", response_class=HTMLResponse)
async def diff(request: Request, output_hash_1_hex: str, output_hash_2_hex: str):

    # Reorder inputs
    # This is to get a deterministic cache key
    output_hash_1_hex, output_hash_2_hex = sorted(
        [output_hash_1_hex, output_hash_2_hex]
    )

    output_hash_1 = codecs.decode(output_hash_1_hex, "hex")  # type: ignore
    output_hash_2 = codecs.decode(output_hash_2_hex, "hex")  # type: ignore

    result1, result2 = await get_derivation_output_results_unique(
        output_hash_1, output_hash_2
    )

    # Uvloop has a nasty bug https://github.com/MagicStack/uvloop/issues/317
    # To work around this we run the fetching/unpacking in a separate blocking thread
    def fetch_unpack_nar(url, location):
        import subprocess

        loc_base = os.path.basename(location)
        loc_dir = os.path.dirname(location)

        try:
            os.mkdir(loc_dir)
        except FileExistsError:
            pass

        with requests.get(url, stream=True) as r:
            r.raise_for_status()
            p = subprocess.Popen(
                ["nix-nar-unpack", loc_base], stdin=subprocess.PIPE, cwd=loc_dir
            )
            for chunk in r.iter_content(chunk_size=512):
                p.stdin.write(chunk)
            p.stdin.close()
            p.wait(timeout=0.5)

        # Ensure correct mtime
        for subl in (
            (os.path.join(dirpath, f) for f in (dirnames + filenames))
            for (dirpath, dirnames, filenames) in os.walk(location)
        ):
            for path in subl:
                os.utime(path, (1, 1))
        os.utime(location, (1, 1))

    async def process_result(result, tmpdir, outbase) -> str:
        # Fetch narinfo
        narinfo = json.loads(
            (await get_combined_rpc().GetValue(api_pb2.ValueRequest(Digest=result.output_hash))).Value  # type: ignore
        )
        nar_hash = narinfo["narHash"].split(":")[-1]

        # Get store prefix
        output = await result.output

        store_base = output.store_path.split("/")[-1]
        store_prefix = store_base.split("-")[0]

        unpack_dir = os.path.join(tmpdir, store_base, outbase)
        nar_url = "/".join((settings.binary_cache_proxy, "nar", store_prefix, nar_hash))

        await asyncio.get_running_loop().run_in_executor(
            None, fetch_unpack_nar, nar_url, unpack_dir
        )

        return unpack_dir

    # TODO: Async tempfile
    with tempfile.TemporaryDirectory(prefix="trustix-ui-dash-diff") as tmpdir:
        dir_a, dir_b = await asyncio.gather(
            process_result(result1, tmpdir, "A"),
            process_result(result2, tmpdir, "B"),
        )

        dir_a_rel = os.path.join(os.path.basename(os.path.dirname(dir_a)), "A")
        dir_b_rel = os.path.join(os.path.basename(os.path.dirname(dir_b)), "B")

        proc = await asyncio.create_subprocess_shell(
            shlex.join(["diffoscope", "--html", "-", dir_a_rel, dir_b_rel]),
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.PIPE,
            cwd=tmpdir,
        )
        stdout, stderr = await proc.communicate()

    # Diffoscope returns non-zero on paths that have a diff
    # Instead use stderr as a heurestic if the call went well or not
    if stderr:
        raise ValueError(stderr)

    return stdout
