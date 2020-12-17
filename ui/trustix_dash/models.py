from tortoise import models
from tortoise import fields


app_name = "trustix_dash"


class Log(models.Model):
    name = fields.CharField(
        max_length=55,
        pk=True,
        index=True,
    )
    tree_size = fields.IntField()

    def __str__(self):
        return self.name


class Evaluation(models.Model):
    commit = fields.CharField(
        max_length=40,
        index=True,
        pk=True,
    )
    timestamp = fields.DatetimeField()  # Commit timestamp

    def __str__(self):
        return self.commit


class Derivation(models.Model):
    drv = fields.CharField(
        max_length=40,
        index=True,
        pk=True,
        unique=True,
    )
    attr = fields.CharField(
        max_length=255,
        index=True,
    )
    system = fields.CharField(
        max_length=255,
        index=True,
    )

    # All build dependencies (recursive)
    refs_all = fields.ManyToManyField(f"{app_name}.Derivation", blank=True)

    # All directy build dependencies (non-recursive)
    refs_direct = fields.ManyToManyField(f"{app_name}.Derivation", blank=True)

    evaluations = fields.ManyToManyField(f"app_name.Evaluation")

    def __str__(self):
        return self.drv


class DerivationOutput(models.Model):
    input_hash = fields.BinaryField(
        max_length=64,
        # index=True,
        # unique=True,
    )
    derivation = fields.ForeignKeyField(f"{app_name}.Derivation", on_delete=fields.CASCADE)
    output = fields.CharField(
        max_length=40,
        index=True,
    )
    store_path = fields.CharField(
        max_length=255,
        index=True,
    )

    def __str__(self):
        return self.store_path

    class Meta:
        unique_together = (("derivation", "output"),)


class DerivationOutputResult(models.Model):
    output = fields.ForeignKeyField(
        f"{app_name}.DerivationOutput",
        to_field="input_hash",
        on_delete=fields.CASCADE,
        db_constraint=False,
    )
    output_hash = fields.BinaryField(max_length=255)
    log = fields.ForeignKeyField(
        f"{app_name}.Log",
        on_delete=fields.CASCADE,
    )

    def __str__(self):
        return f"{self.log_id}({self.output.derivation_id, self.output.output})"

    class Meta:
        unique_together = (("output", "log"),)
