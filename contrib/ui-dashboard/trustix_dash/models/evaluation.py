from tortoise import models
from tortoise import fields

from trustix_dash.models.meta import BaseMeta


class Evaluation(models.Model):
    commit = fields.CharField(
        max_length=40,
        index=True,
        pk=True,
    )

    # TODO: Reliably get a timestamp of eval (from hydra api?)
    # timestamp = fields.DatetimeField()  # Commit timestamp

    def __str__(self):
        return self.commit

    class Meta(BaseMeta):
        pass
