import pathlib

from litestar import Litestar, get
from litestar.response import Template, Response
from litestar.template import TemplateConfig
from litestar.contrib.jinja import JinjaTemplateEngine

from app.utils import (
    read_snippets,
    get_runtime_metadata,
)


SNIPPETS = read_snippets()
RUNTIME_METADATA = get_runtime_metadata()


@get("/")
async def index() -> Template:
    return Template(
        "home.html.jinja",
        context={
            "items": [(bit, name) for bit, (name, _) in SNIPPETS.items()],
            "repo": {
                "branch": RUNTIME_METADATA.branch,
                "commit": RUNTIME_METADATA.commit,
                "origin": RUNTIME_METADATA.origin,
            },
        },
    )


@get("/snippets/{mask:int}")
async def snippets(mask: int) -> Response:
    css = "\n".join(css for bit, (_, css) in SNIPPETS.items() if mask & (1 << bit))

    return Response(
        content=css.strip(),
        media_type="text/css",
    )


app = Litestar(
    route_handlers=[index, snippets],
    openapi_config=None,
    template_config=TemplateConfig(
        directory=pathlib.Path("templates"),
        engine=JinjaTemplateEngine,
    ),
)
