from trustix_dash import index_logs, index_eval, get_derivation_outputs
from tortoise import run_async
from tortoise import Tortoise


async def init():

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

    # # TODO: Remove and use aerich instead
    # await Tortoise.generate_schemas(safe=True)

    # import pdb, traceback, sys
    # try:
    #     obj = await get_derivation_outputs("s6rn4jz1sin56rf4qj5b5v8jxjm32hlk-hello-2.10.drv")
    # except:
    #     extype, value, tb = sys.exc_info()
    #     traceback.print_exc()
    #     pdb.post_mortem(tb)
    # else:
    #     print(obj)

    # print(
    #     await get_derivation_outputs("s6rn4jz1sin56rf4qj5b5v8jxjm32hlk-hello-2.10.drv")
    # )

    # await index_logs()

    commit_sha = "e9158eca70ae59e73fae23be5d13d3fa0cfc78b4"
    await index_eval(commit_sha)


run_async(init())
