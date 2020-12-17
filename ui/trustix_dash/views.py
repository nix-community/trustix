from dash_main.api import get_derivation_outputs
from django.http import HttpResponse
from django.template import loader


def index(request):
    template = loader.get_template("index.html")
    context = {}

    d = "s6rn4jz1sin56rf4qj5b5v8jxjm32hlk-hello-2.10.drv"
    for result in get_derivation_outputs(d):
        output = result.output
        print(output)

    return HttpResponse(template.render(context, request))
