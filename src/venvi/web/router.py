from pathlib import Path
from typing import Sequence

from fastapi import APIRouter, Request, Depends
from fastapi.templating import Jinja2Templates
from sqlalchemy.ext.asyncio import AsyncSession
from sqlmodel import select

from venvi.core.db import get_session
from venvi.models.hackathon import Hackathon

router = APIRouter()

# Path to the templates directory: src/venvi/templates
# This file is in: src/venvi/web/router.py
BASE_DIR = Path(__file__).resolve().parent.parent
templates = Jinja2Templates(directory=str(BASE_DIR / "templates"))

@router.get("/", include_in_schema=False)
async def index(request: Request):
    return templates.TemplateResponse(request, "index.html")

@router.get("/partials/hackathons", include_in_schema=False)
async def get_hackathons_partial(
    request: Request, 
    session: AsyncSession = Depends(get_session)
):
    query = select(Hackathon).order_by(Hackathon.date_start)
    result = await session.execute(query)
    hackathons: Sequence[Hackathon] = result.scalars().all()
    
    return templates.TemplateResponse(
        request, 
        "partials/hackathon_list.html", 
        {"hackathons": hackathons}
    )
