from django.contrib import admin
from .models import (
    Log,
    Evaluation,
    Derivation,
    DerivationOutput,
    DerivationOutputResult,
)


class LogAdmin(admin.ModelAdmin):
    pass


class EvaluationAdmin(admin.ModelAdmin):
    pass


class DerivationAdmin(admin.ModelAdmin):
    raw_id_fields = ("evaluations",)


class DerivationOutputAdmin(admin.ModelAdmin):
    raw_id_fields = ("derivation",)


class DerivationOutputResultAdmin(admin.ModelAdmin):
    raw_id_fields = (
        "output",
        "log",
    )


admin.site.register(Log, LogAdmin)
admin.site.register(Evaluation, EvaluationAdmin)
admin.site.register(Derivation, DerivationAdmin)
admin.site.register(DerivationOutput, DerivationOutputAdmin)
admin.site.register(DerivationOutputResult, DerivationOutputResultAdmin)
