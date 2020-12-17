from django.db import models


class Log(models.Model):
    name = models.CharField(
        max_length=55,
        primary_key=True,
        db_index=True,
    )
    tree_size = models.IntegerField()

    def __str__(self):
        return self.name


class Evaluation(models.Model):
    commit = models.CharField(
        max_length=40,
        db_index=True,
        primary_key=True,
    )
    timestamp = models.DateTimeField()  # Commit timestamp

    def __str__(self):
        return self.commit


class Derivation(models.Model):
    drv = models.CharField(
        max_length=40,
        db_index=True,
        primary_key=True,
    )
    attr = models.CharField(
        max_length=255,
        db_index=True,
    )
    system = models.CharField(
        max_length=255,
        db_index=True,
    )

    # All build dependencies (recursive)
    refs_all = models.ManyToManyField("self", blank=True)

    # All directy build dependencies (non-recursive)
    refs_direct = models.ManyToManyField("self", blank=True)

    evaluations = models.ManyToManyField(Evaluation)

    def __str__(self):
        return self.drv


class DerivationOutput(models.Model):
    input_hash = models.BinaryField(
        max_length=64,
        db_index=True,
        unique=True,
    )
    derivation = models.ForeignKey(Derivation, on_delete=models.CASCADE)
    output = models.CharField(
        max_length=40,
        db_index=True,
    )
    store_path = models.CharField(
        max_length=255,
        db_index=True,
    )

    def __str__(self):
        return self.store_path

    class Meta:
        constraints = [
            models.UniqueConstraint(
                fields=["derivation", "output"], name="unique_output"
            ),
        ]


class DerivationOutputResult(models.Model):
    output = models.ForeignKey(
        DerivationOutput,
        to_field="input_hash",
        on_delete=models.CASCADE,
        db_constraint=False,
    )
    output_hash = models.BinaryField(max_length=255)
    log = models.ForeignKey(
        Log,
        on_delete=models.CASCADE,
    )

    def __str__(self):
        return f"{self.log_id}({self.output.derivation_id, self.output.output})"

    class Meta:
        constraints = [
            models.UniqueConstraint(
                fields=["output", "log"], name="unique_output_result"
            ),
        ]
