from django.core.management.base import BaseCommand
from dash_main.api import check_reprod


class Command(BaseCommand):
    def handle(self, *args, **options):
        check_reprod("hello")
