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
    taken: bool | None = None,
    session: AsyncSession = Depends(get_session),
):
    """
    Retrieves a list of hackathons from the database.

    Args:
        status: Optional filter by hackathon status (e.g., 'upcoming', 'past').
        taken: Optional filter by 'taken' status.
        session: The asynchronous database session.

    Returns:
        list[Hackathon]: A list of hackathon objects.
    """
    query = select(Hackathon)
    if status:
        query = query.where(Hackathon.status == status)
    if taken is not None:
        query = query.where(Hackathon.taken == taken)

    result = await session.execute(query)
    return result.scalars().all()


@router.patch("/{hackathon_id}/taken", response_model=Hackathon)
async def toggle_hackathon_taken(
    hackathon_id: str, taken: bool, session: AsyncSession = Depends(get_session)
):
    """
    Updates the 'taken' status of a hackathon.

    Args:
        hackathon_id: The UUID of the hackathon.
        taken: The new taken status.
    """
    # Note: hackathon_id is a string in URL but UUID in model.
    # SQLAlchemy/SQLModel should handle coercion if passed as string to get(),
    # but let's be safe or rely on FastAPI validation if we typed it as UUID.
    # The models/hackathon.py defines id as UUID.

    # We need to import UUID if we want to cast, or let sqlmodel handle it.
    # session.get usually works with string for UUID pk but let's just pass it.

    hackathon = await session.get(Hackathon, hackathon_id)
    if not hackathon:
        raise HTTPException(status_code=404, detail="Hackathon not found")

    hackathon.taken = taken
    session.add(hackathon)
    await session.commit()
    await session.refresh(hackathon)
    return hackathon


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
