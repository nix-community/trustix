# from fastapi.staticfiles import StaticFiles
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
# from pydantic import (
#     BaseModel,
# )
import os.path

from trustix_dash.api import (
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


@app.get("/drv/{store_path}", response_class=HTMLResponse)
async def drv(request: Request):
    ctx = await make_context(request)
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
