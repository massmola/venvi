"""
API endpoints for managing and viewing events from multiple sources.
"""

from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.ext.asyncio import AsyncSession
from sqlmodel import select

from venvi.core.db import get_session
from venvi.models.odh_event import ODHEvent
from venvi.services.odh import sync_odh_events

router = APIRouter(prefix="/events", tags=["events"])


@router.get("/odh", response_model=list[ODHEvent])
async def read_odh_events(session: AsyncSession = Depends(get_session)):
    """
    Retrieves a list of South Tyrol Open Data Hub events from the database.

    Args:
        session: The asynchronous database session.

    Returns:
        list[ODHEvent]: A list of ODH event objects.
    """
    query = select(ODHEvent)
    result = await session.execute(query)
    return result.scalars().all()


@router.post("/odh/sync")
async def sync_odh_data(session: AsyncSession = Depends(get_session)):
    """
    Triggers an on-demand synchronization with South Tyrol Open Data Hub.

    Args:
        session: The asynchronous database session.

    Returns:
        dict: A status message and the count of new items added.

    Raises:
        HTTPException: If the synchronization process fails.
    """
    try:
        count = await sync_odh_events(session)
        return {"message": "Sync complete", "new_items": count}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e)) from e
