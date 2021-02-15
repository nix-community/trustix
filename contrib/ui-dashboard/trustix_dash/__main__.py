from trustix_dash import (
    index_logs,
    index_eval,
    get_derivation_outputs,
)
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

    # TODO: Remove and use aerich instead (blocked by https://github.com/tortoise/aerich/issues/63 )
    await Tortoise.generate_schemas(safe=True)

    await index_logs()

    commit_sha = "fb6f9b7eb0aa8629776ea32d2be6eaf660a22535"
    await index_eval(commit_sha)

    import pdb, traceback, sys
    try:
        obj = await get_derivation_outputs("/nix/store/r4vd7wq0j9py6vrzsv7r4ggpcg95jkys-hello-2.10.drv")
    except:
        extype, value, tb = sys.exc_info()
        traceback.print_exc()
        pdb.post_mortem(tb)
    else:
        print(obj)


if __name__ == '__main__':
    run_async(init())
