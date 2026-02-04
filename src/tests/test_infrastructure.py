from pathlib import Path

import pytest
from httpx import AsyncClient


def test_static_directories_exist() -> None:
    """Verify that all required static subdirectories exist."""
    base_dir = Path(__file__).resolve().parent.parent / "venvi" / "static"
    subdirs = ["css", "img", "js"]

    assert base_dir.exists(), f"Static directory {base_dir} does not exist"
    for subdir in subdirs:
        path = base_dir / subdir
        assert (
            path.is_dir()
        ), f"Required subdirectory {path} is missing or not a directory"


@pytest.mark.asyncio
async def test_static_file_serving(client: AsyncClient) -> None:
    """Verify that the application serves static files correctly."""
    # Test serving the placeholder CSS file
    response = await client.get("/static/css/style.css")
    assert response.status_code == 200
    assert "text/css" in response.headers["content-type"]
    assert "background-color" in response.text
