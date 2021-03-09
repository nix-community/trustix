from tortoise import models
from tortoise import fields

from trustix_dash.models.meta import BaseMeta


class Evaluation(models.Model):
    commit = fields.CharField(
        max_length=40,
        index=True,
        pk=True,
    )

    timestamp = fields.DatetimeField(auto_now_add=True)

    def __str__(self):
        return self.commit

    class Meta(BaseMeta):
        pass
