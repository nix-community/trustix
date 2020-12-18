from trustix_dash import fields as trustix_fields
from tortoise import models
from tortoise import fields


app_name = "trustix_dash"


class BaseMeta:
    app = app_name


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


class Derivation(models.Model):
    drv = fields.CharField(
        max_length=40,
        index=True,
        pk=True,
        unique=True,
    )
    system = fields.CharField(
        max_length=255,
        index=True,
    )

    # All build dependencies (recursive)
    refs_all = fields.ManyToManyField(
        f"{app_name}.Derivation", related_name="refsall", through="derivation_refs_all"
    )

    # # All directy build dependencies (non-recursive)
    refs_direct = fields.ManyToManyField(
        f"{app_name}.Derivation",
        related_name="refsdirect",
        through="derivation_refs_direct",
    )

    evaluations = fields.ManyToManyField(f"{app_name}.Evaluation")

    def __str__(self):
        return self.drv

    class Meta(BaseMeta):
        pass


class DerivationAttr(models.Model):
    derivation = fields.ForeignKeyField(
        f"{app_name}.Derivation", on_delete=fields.CASCADE
    )
    attr = fields.CharField(
        max_length=255,
        index=True,
    )

    def __str__(self):
        return self.attr

    class Meta(BaseMeta):
        unique_together = (("derivation", "attr"),)


class DerivationOutput(models.Model):
    # Input hash == store path prefix
    input_hash = trustix_fields.BinaryField(
        max_length=25,
        pk=True,
        index=True,
    )
    # TODO: Make ManyToMany (2 different drvs with the same src but different fetchers are distinct drvs with the same outputs)
    derivation = fields.ForeignKeyField(
        f"{app_name}.Derivation", on_delete=fields.CASCADE
    )
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

    class Meta(BaseMeta):
        unique_together = (("derivation", "output"),)


class DerivationOutputResult(models.Model):
    output = fields.ForeignKeyField(
        f"{app_name}.DerivationOutput",
        to_field="input_hash",
        on_delete=fields.CASCADE,
        db_constraint=False,
    )
    output_hash = trustix_fields.BinaryField(max_length=40)
    log = fields.ForeignKeyField(
        f"{app_name}.Log",
        on_delete=fields.CASCADE,
    )

    def __str__(self):
        return f"{self.log_id}({self.output.derivation_id, self.output.output})"

    class Meta(BaseMeta):
        unique_together = (("output", "log"),)
