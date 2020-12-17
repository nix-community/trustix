from django.core.management.base import BaseCommand
from dash_main.api import index_eval_outputs


class Command(BaseCommand):
    def handle(self, *args, **options):
        commit_sha = "e9158eca70ae59e73fae23be5d13d3fa0cfc78b4"
        index_eval_outputs(commit_sha)
