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
    Optional,
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
import os.path

from trustix_nix_reprod.staticfiles import StaticFiles

from pydantic import BaseModel as BaseModel

from trustix_nix_reprod.api import (
    get_derivation_reproducibility,
    get_attrs_reproducibility,
    search_derivations,
    suggest_attrs,
)
from trustix_nix_reprod.diff import json_diff
from trustix_nix_reprod.api import diff as diff_api
from trustix_nix_reprod.api.models import (
    DerivationReproducibility,
    SuggestResponse,
    SearchResponse,
    AttrsResponse,
    DiffResponse,
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
templates.env.globals["diffoscope_render"] = template_lib.diffoscope_render


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
    extra: Optional[Dict] = None,
) -> Dict:

    ctx = {
        "request": request,
        "title": "Trustix R13Y" + (" ".join((" - ", title)) if title else ""),
        "drv_placeholder": settings.placeholder_attr,
        "model": model,
    }

    if extra:
        ctx.update(extra)

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


@app.get("/diff/{output_hash_1_hex}/{output_hash_2_hex}", response_model=DiffResponse)
async def diff(request: Request, output_hash_1_hex: str, output_hash_2_hex: str):
    # Reorder inputs
    # This is to get a deterministic cache key
    output_hash_1_hex, output_hash_2_hex = sorted(
        [output_hash_1_hex, output_hash_2_hex]
    )
    data = await diff_api(output_hash_1_hex, output_hash_2_hex)

    def diff_template_context():
        return make_context(request, data, extra={
            "narinfo_diff": json_diff(
                data.narinfo[output_hash_1_hex],
                data.narinfo[output_hash_2_hex],
            )
        })

    return render_model(
        request,
        "diff.jinja2",
        diff_template_context,
        data,
    )
