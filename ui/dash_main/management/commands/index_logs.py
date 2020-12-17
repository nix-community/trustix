from django.core.management.base import BaseCommand
from dash_main.api import index_logs


class Command(BaseCommand):
    def handle(self, *args, **options):
        index_logs()
