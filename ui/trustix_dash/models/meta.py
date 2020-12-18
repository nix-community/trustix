from tortoise import models
from tortoise import fields


app_name = "trustix_dash"


class BaseMeta:
    app = app_name
