from tortoise.contrib.fastapi import register_tortoise
from fastapi.templating import (
    Jinja2Templates,
)
from fastapi.responses import (
    HTMLResponse,
)
from fastapi import (
    FastAPI,
    Request,
)
from typing import (
    Optional,
    Dict,
    List,
)
from trustix_dash.models import (
    DerivationAttr,
)
from tortoise import Tortoise
import urllib.parse
import asyncio
import os.path
import shlex

from trustix_dash.api import (
    get_derivation_output_results,
    get_derivation_outputs,
    evaluation_list,
)


templates = Jinja2Templates(
    directory=os.path.join(os.path.dirname(__file__), "templates")
)


app = FastAPI()


@app.on_event("startup")
async def startup_event():
    await Tortoise.init(
        {
            "connections": {
                "default": "sqlite://db.sqlite3",
            },
            "apps": {
                "trustix_dash": {
                    "models": ["trustix_dash.models"],
                }
            },
            "use_tz": False,
            "timezone": "UTC",
        }
    )


@app.on_event("shutdown")
async def close_orm():
    await Tortoise.close_connections()


async def make_context(
    request: Request,
    title: str = "",
    selected_evaluation: Optional[str] = None,
    extra: Optional[Dict] = None,
) -> Dict:

    evaluations = await evaluation_list()
    if selected_evaluation and selected_evaluation not in evaluations:
        evaluations.insert(0, selected_evaluation)

    if not selected_evaluation:
        try:
            selected_evaluation = evaluations[0]
        except IndexError:
            pass

    ctx = {
        "evaluations": evaluations,
        "selected_evaluation": selected_evaluation,
        "request": request,
        "title": "Trustix R13Y" + (" ".join((" - ", title)) if title else ""),
        "drv_placeholder": "hello.x86_64-linux",
    }

    if extra:
        ctx.update(extra)

    return ctx


@app.get("/", response_class=HTMLResponse)
async def index(request: Request):
    ctx = await make_context(request)
    return templates.TemplateResponse("index.jinja2", ctx)


@app.get("/drv/{drv_path}", response_class=HTMLResponse)
async def drv(request: Request, drv_path: str):

    drv_path = urllib.parse.unquote(drv_path)
    drvs = await get_derivation_outputs(drv_path)

    unreproduced_paths: Dict[str, List[str]] = {}
    reproduced_paths: Dict[str, List[str]] = {}
    missing_paths: Dict[str, List[str]] = {}  # Paths not built by any known log

    for drv in drvs:
        for output in drv.derivationoutputs:  # type: ignore
            output_hashes = set(
                result.output_hash for result in output.derivationoutputresults
            )

            if output_hashes:
                print([result.id for result in output.derivationoutputresults])

            if not output_hashes:
                missing_paths.setdefault(drv.drv, []).append(output.output)  # type: ignore

            elif len(output_hashes) == 1:
                reproduced_paths.setdefault(drv.drv, []).append(output.output)  # type: ignore

            elif len(output_hashes) > 1:
                unreproduced_paths.setdefault(drv.drv, []).append(output.output)  # type: ignore

            else:
                raise RuntimeError("Logic error")

    ctx = await make_context(
        request,
        extra={
            "unreproduced_paths": unreproduced_paths,
            "reproduced_paths": reproduced_paths,
            "missing_paths": missing_paths,
        },
    )

    return templates.TemplateResponse("drv.jinja2", ctx)


# @app.post("/search")
# async def search(request: Request, attr: Optional[str] = None):
#     return {}


@app.get("/suggest/{attr}", response_model=List[str])
async def suggest(request: Request, attr: str):
    if len(attr) < 3:
        raise ValueError("Prefix too short")

    resp = await DerivationAttr.filter(attr__startswith=attr).only("attr")
    return [drv_attr.attr for drv_attr in resp]


@app.get("/diff/{output1}/{output2}", response_class=HTMLResponse)
async def diff(request: Request, output1: int, output2: int):
    result1, result2 = await get_derivation_output_results(output1, output2)

    store_path_1 = (await result1.output).store_path
    store_path_2 = (await result2.output).store_path

    proc = await asyncio.create_subprocess_shell(
        shlex.join(["diffoscope", "--html", "-", store_path_1, store_path_2]),
        stdout=asyncio.subprocess.PIPE,
        stderr=asyncio.subprocess.PIPE,
    )

    stdout, stderr = await proc.communicate()

    # Diffoscope returns non-zero on paths being different
    # Instead use stderr as a heurestic if the call went well or not
    if stderr:
        raise ValueError(stderr)

    return stdout
