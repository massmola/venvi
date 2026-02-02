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
@router.get("/odh", response_model=list[ODHEvent])
async def read_odh_events(
    taken: bool | None = None, session: AsyncSession = Depends(get_session)
):
    """
    Retrieves a list of South Tyrol Open Data Hub events from the database.

    Args:
        taken: Optional filter by 'taken' status.
        session: The asynchronous database session.

    Returns:
        list[ODHEvent]: A list of ODH event objects.
    """
    query = select(ODHEvent)
    if taken is not None:
        query = query.where(ODHEvent.taken == taken)

    result = await session.execute(query)
    return result.scalars().all()


@router.patch("/odh/{event_id}/taken", response_model=ODHEvent)
async def toggle_odh_event_taken(
    event_id: str, taken: bool, session: AsyncSession = Depends(get_session)
):
    """
    Updates the 'taken' status of an ODH event.

    Args:
        event_id: The ID of the event.
        taken: The new taken status.
        session: The asynchronous database session.

    Returns:
        ODHEvent: The updated event.

    Raises:
        HTTPException: If the event is not found.
    """
    event = await session.get(ODHEvent, event_id)
    if not event:
        raise HTTPException(status_code=404, detail="Event not found")

    event.taken = taken
    session.add(event)
    await session.commit()
    await session.refresh(event)
    return event


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
