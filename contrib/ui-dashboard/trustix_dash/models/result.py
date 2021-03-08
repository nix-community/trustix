from tortoise import models
from tortoise import fields

from trustix_dash.models import fields as trustix_fields
from trustix_dash.models.meta import BaseMeta, app_name


class DerivationOutputResult(models.Model):
    output = fields.ForeignKeyField(
        f"{app_name}.DerivationOutput",
        to_field="input_hash",
        on_delete=fields.CASCADE,
        db_constraint=False,
        index=True,
    )
    # TODO: Turn into indexed CharField
    output_hash = trustix_fields.BinaryField(max_length=40)
    log = fields.ForeignKeyField(
        f"{app_name}.Log",
        on_delete=fields.CASCADE,
    )

    # def __str__(self):
    #     return f"{self.output_id}"
    #     # return f"{self.log_id}({self.output.derivation_id, self.output.output})"

    class Meta(BaseMeta):
        unique_together = (("output", "log"),)
