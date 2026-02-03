"""
API endpoints for managing and viewing hackathon data.
"""

from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.ext.asyncio import AsyncSession
from sqlmodel import select

from venvi.core.db import get_session
from venvi.models.hackathon import Hackathon
from venvi.services.ingestion import sync_hackathons

router = APIRouter(prefix="/hackathons", tags=["hackathons"])


@router.get("/", response_model=list[Hackathon])
async def read_hackathons(
    status: str | None = None,
    session: AsyncSession = Depends(get_session),
):
    """
    Retrieves a list of hackathons from the database.

    Args:
        status: Optional filter by hackathon status (e.g., 'upcoming', 'past').
        session: The asynchronous database session.

    Returns:
        list[Hackathon]: A list of hackathon objects.
    """
    query = select(Hackathon)
    if status:
        query = query.where(Hackathon.status == status)

    result = await session.execute(query)
    return result.scalars().all()


@router.post("/sync")
async def sync_data(session: AsyncSession = Depends(get_session)):
    """
    Triggers an on-demand synchronization with external data sources.

    Args:
        session: The asynchronous database session.

    Returns:
        dict: A status message and the count of new items added.

    Raises:
        HTTPException: If the synchronization process fails.
    """
    try:
        count = await sync_hackathons(session)
        return {"message": "Sync complete", "new_items": count}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e)) from e
