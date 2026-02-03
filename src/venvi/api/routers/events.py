from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.ext.asyncio import AsyncSession
from sqlmodel import select

from venvi.core.db import get_session
from venvi.models.event import Event
from venvi.services.ingestion import sync_all_events

router = APIRouter(prefix="/events", tags=["events"])


@router.get("/", response_model=list[Event])
async def read_events(
    category: str | None = None,
    source: str | None = None,
    session: AsyncSession = Depends(get_session),
):
    """
    Retrieves a list of events from the database.
    """
    query = select(Event).order_by(Event.date_start.asc())
    if category:
        query = query.where(Event.category == category)
    if source:
        query = query.where(Event.source_name == source)

    result = await session.execute(query)
    return result.scalars().all()


@router.post("/sync")
async def sync_events(session: AsyncSession = Depends(get_session)):
    """
    Triggers an on-demand synchronization with all event providers.
    """
    try:
        results = await sync_all_events(session)
        total_new = sum(results.values())
        return {"message": "Sync complete", "new_items": results, "total": total_new}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e)) from e
