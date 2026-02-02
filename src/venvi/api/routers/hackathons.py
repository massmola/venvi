from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.ext.asyncio import AsyncSession
from sqlmodel import select
from venvi.core.db import get_session
from venvi.models.hackathon import Hackathon
from venvi.services.ingestion import sync_hackathons

router = APIRouter(prefix="/hackathons", tags=["hackathons"])


@router.get("/", response_model=list[Hackathon])
async def read_hackathons(
    status: str | None = None, session: AsyncSession = Depends(get_session)
):
    query = select(Hackathon)
    if status:
        query = query.where(Hackathon.status == status)

    result = await session.execute(query)
    return result.scalars().all()


@router.post("/sync")
async def sync_data(session: AsyncSession = Depends(get_session)):
    try:
        count = await sync_hackathons(session)
        return {"message": "Sync complete", "new_items": count}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e)) from e
