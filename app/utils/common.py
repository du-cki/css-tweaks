import pygit2
from pathlib import Path

from typing import Dict, Tuple, NamedTuple

SNIPPETS_DIR = Path("snippets")


def read_snippets() -> Dict[int, Tuple[str, str]]:
    css = {}

    snippets = SNIPPETS_DIR.glob("*.css")
    for snippet in snippets:
        with open(snippet, "r") as f:
            id, name = snippet.stem.split("-", 1)
            parsed_name = " ".join(name.split("-")).title()

            css[int(id)] = (parsed_name, f"/* {parsed_name} */\n{f.read()}")

    return css


class RuntimeMetadata(NamedTuple):
    branch: str
    commit: str
    origin: str


def get_runtime_metadata() -> RuntimeMetadata:
    repo = pygit2.Repository(".git")
    branch = repo.head.shorthand
    commit = repo.head.peel(pygit2.Commit)  # pyright: ignore[reportArgumentType,reportCallIssue]
    origin = repo.remotes["origin"].url

    if not origin:
        raise ValueError("No remote origin found in the current repository.")

    return RuntimeMetadata(
        branch=branch,
        commit=commit.id.hex,
        origin=origin,
    )


__all__ = (
    "read_snippets",
    "get_runtime_metadata",
)
