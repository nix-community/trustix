from django.core.management.base import BaseCommand
from dash_main.api import check_reprod


class Command(BaseCommand):
    def handle(self, *args, **options):
        d = "s6rn4jz1sin56rf4qj5b5v8jxjm32hlk-hello-2.10.drv"
        for x in check_reprod(d):
            print(x)
