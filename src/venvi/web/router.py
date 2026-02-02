"""
Web routes for rendering the application's HTML interface.
"""

from collections.abc import Sequence
from pathlib import Path

from fastapi import APIRouter, Depends, Request, Response
from fastapi.templating import Jinja2Templates
from sqlalchemy.ext.asyncio import AsyncSession
from sqlmodel import select

from venvi.core.db import get_session
from venvi.models.hackathon import Hackathon

router = APIRouter()

# Path to the templates directory: src/venvi/templates
BASE_DIR = Path(__file__).resolve().parent.parent
templates = Jinja2Templates(directory=str(BASE_DIR / "templates"))


@router.get("/", include_in_schema=False)
async def index(request: Request) -> Response:
    """
    Renders the main application homepage.

    Args:
        request: The Starlette request object.

    Returns:
        TemplateResponse: The rendered index.html template.
    """
    return templates.TemplateResponse(request, "index.html")


@router.get("/partials/hackathons", include_in_schema=False)
async def get_hackathons_partial(
    request: Request, session: AsyncSession = Depends(get_session)
) -> Response:
    """
    Renders the hackathon list partial for dynamic HTMX updates.

    Args:
        request: The Starlette request object.
        session: The asynchronous database session.

    Returns:
        TemplateResponse: The rendered partials/hackathon_list.html template.
    """
    query = select(Hackathon).order_by(Hackathon.date_start.asc())
    result = await session.execute(query)
    hackathons: Sequence[Hackathon] = result.scalars().all()

    return templates.TemplateResponse(
        request, "partials/hackathon_list.html", {"hackathons": hackathons}
    )
