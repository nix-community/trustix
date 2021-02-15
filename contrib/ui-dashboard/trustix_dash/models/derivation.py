from tortoise import models
from tortoise import fields

from trustix_dash.models import fields as trustix_fields
from trustix_dash.models.meta import BaseMeta, app_name


__all__ = (
    "DerivationRefRecursive",
    "DerivationRefDirect",
    "DerivationOutput",
    "DerivationAttr",
    "DerivationEval",
    "Derivation",
)


class Derivation(models.Model):
    drv = fields.CharField(
        max_length=120,
        index=True,
        pk=True,
        unique=True,
    )
    system = fields.CharField(
        max_length=255,
        index=True,
    )

    def __str__(self):
        return self.drv

    class Meta(BaseMeta):
        pass


class DerivationEval(models.Model):

    drv = fields.ForeignKeyField(
        f"{app_name}.Derivation",
        on_delete=fields.CASCADE,
        index=True,
    )
    eval = fields.ForeignKeyField(
        f"{app_name}.Evaluation",
        on_delete=fields.CASCADE,
        index=True,
    )

    def __str__(self):
        return "->".join((self.eval_id, self.drv_id))

    class Meta(BaseMeta):
        pass


def MkAbstractDerivationRef(name):
    class AbstractDerivationRef(models.Model):
        drv = fields.ForeignKeyField(
            f"{app_name}.Derivation",
            on_delete=fields.CASCADE,
            related_name=f"from_ref_{name}",
            index=True,
        )
        referrer = fields.ForeignKeyField(
            f"{app_name}.Derivation",
            on_delete=fields.CASCADE,
            related_name=f"to_ref_{name}",
            index=True,
        )

        class Meta(BaseMeta):
            abstract = True

    return AbstractDerivationRef


class DerivationRefRecursive(MkAbstractDerivationRef("recursive")):  # type: ignore
    pass


class DerivationRefDirect(MkAbstractDerivationRef("direct")):  # type: ignore
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
        max_length=120,
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
