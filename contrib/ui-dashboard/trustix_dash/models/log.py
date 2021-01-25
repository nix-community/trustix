from tortoise import models
from tortoise import fields

from trustix_dash.models.meta import BaseMeta


class Log(models.Model):
    name = fields.CharField(
        max_length=55,
        index=True,
    )
    tree_size = fields.IntField()

    def __str__(self):
        return self.name

    class Meta(BaseMeta):
        pass
