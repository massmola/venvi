"""
Services for aggregating and synchronizing hackathon data from external providers.
"""

from typing import Any

import httpx
from sqlalchemy.ext.asyncio import AsyncSession
from sqlmodel import select

from venvi.models.hackathon import Hackathon

EURO_HACKATHONS_API = "https://euro-hackathons.vercel.app/api/hackathons"




async def fetch_euro_hackathons() -> list[dict[str, Any]]:
    """
    Fetches upcoming hackathons from the Euro Hackathons API.

    Returns:
        list[dict]: A list of raw hackathon data dictionaries.

    Raises:
        httpx.HTTPStatusError: If the API request fails.
    """
    async with httpx.AsyncClient() as client:
        response = await client.get(EURO_HACKATHONS_API, params={"status": "upcoming"})
        response.raise_for_status()
        return response.json().get("data", [])


async def sync_hackathons(session: AsyncSession) -> int:
    """
    Synchronizes the local database with the Euro Hackathons API.

    New items are created, and existing items are updated if their data has changed.

    Args:
        session: The asynchronous database session.

    Returns:
        int: The number of new hackathons added.
    """
    raw_hackathons = await fetch_euro_hackathons()
    count = 0

    for raw in raw_hackathons:
        # Check if already exists
        statement = select(Hackathon).where(Hackathon.id == raw["id"])
        result = await session.execute(statement)
        existing = result.scalar_one_or_none()

        if existing:
            # Update fields if necessary
            updated_instance = Hackathon.model_validate(raw)
            update_data = updated_instance.model_dump(exclude_unset=True)
            for key, value in update_data.items():
                if hasattr(existing, key):
                    setattr(existing, key, value)
            session.add(existing)
        else:
            # Create new
            hackathon = Hackathon.model_validate(raw)
            session.add(hackathon)
            count += 1

    await session.commit()
    return count
