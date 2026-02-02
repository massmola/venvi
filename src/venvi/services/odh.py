"""
Services for fetching and synchronizing events from the South Tyrol Open Data Hub.
"""

from datetime import datetime

import httpx
from sqlalchemy.ext.asyncio import AsyncSession
from sqlmodel import select

from venvi.models.odh_event import ODHEvent

# Base URL for ODH Tourism API
ODH_API_URL = "https://tourism.opendatahub.com/v1/Event"


async def fetch_odh_events(page_size: int = 20) -> list[dict]:
    """
    Fetches events from the Open Data Hub API.

    Args:
        page_size: Number of events to fetch.

    Returns:
        list[dict]: A list of raw event data dictionaries.

    Raises:
        httpx.HTTPStatusError: If the API request fails.
    """
    params = {
        "pagenumber": 1,
        "pagesize": page_size,
    }

    # We might want to filter by date or other parameters in the future.
    # The current goal is basic integration.

    async with httpx.AsyncClient() as client:
        response = await client.get(ODH_API_URL, params=params)
        response.raise_for_status()
        data = response.json()
        return data.get("Items", [])


def map_odh_event(raw: dict) -> ODHEvent:
    """
    Maps raw ODH event data to the ODHEvent model.

    Args:
        raw: Dictionary containing raw event data from ODH.

    Returns:
        ODHEvent: The mapped model instance.
    """
    # Defensive programming for handling optional fields and language selection (prefer 'en', fallback to 'it'/'de')

    def get_start_date(raw_event: dict) -> datetime:
        # Assuming EventDate is a list and we take the first one or the main DateBegin
        # The structure is a bit complex, let's try to find the best date
        # Check 'DateBegin' first
        if raw_event.get("DateBegin"):
            return datetime.fromisoformat(raw_event["DateBegin"])
        return datetime.now()  # Fallback, should not happen for valid events

    def get_end_date(raw_event: dict) -> datetime:
        if raw_event.get("DateEnd"):
            return datetime.fromisoformat(raw_event["DateEnd"])
        return get_start_date(raw_event)

    def get_localized_string(obj: dict | None, key: str) -> str | None:
        if not obj:
            return None

        # Try to find the key in language dictionaries
        for lang in ["en", "it", "de"]:
            lang_data = obj.get(lang)
            if isinstance(lang_data, dict):
                val = lang_data.get(key)
                if val:
                    return val
        return None

    details = raw.get("Detail", {})

    title = get_localized_string(details, "Title") or "Untitled Event"
    # Description is often in 'BaseText' or 'IntroText'
    description = get_localized_string(details, "BaseText") or get_localized_string(
        details, "IntroText"
    )

    # Location
    location_info = raw.get("ContactInfos", {})
    # Just take city from contact info if available, or LocationInfo
    # The structure saw earlier had 'ContactInfos' with 'en', 'it', etc.
    contact_en = location_info.get("en", {})
    city = contact_en.get("City") or "Unknown Location"

    # Image
    image_url = None
    gallery = raw.get("ImageGallery", [])
    if gallery:
        image_url = gallery[0].get("ImageUrl")

    # Try to resolve or generate an ID
    raw_id = raw.get("Id")

    raw_id = raw.get("Id")
    if not raw_id:
        # try to find it in mapping
        mapping = raw.get("Mapping", {})
        for source in mapping.values():
            if "rid" in source:
                raw_id = source["rid"]
                break

    if not raw_id:
        # Last resort fallback
        raw_id = str(hash(title + str(raw.get("DateBegin"))))

    return ODHEvent(
        id=raw_id,
        title=title,
        description=description,
        date_start=get_start_date(raw),
        date_end=get_end_date(raw),
        location=city,
        image_url=image_url,
        source_url=None,  # Detail link might be constructed
        is_new=True,
    )


async def sync_odh_events(session: AsyncSession) -> int:
    """
    Synchronizes the local database with the Open Data Hub API.

    Args:
        session: The asynchronous database session.

    Returns:
        int: The number of new events added.
    """
    raw_events = await fetch_odh_events()
    count = 0

    for raw in raw_events:
        event = map_odh_event(raw)

        # Check if exists
        statement = select(ODHEvent).where(ODHEvent.id == event.id)
        result = await session.execute(statement)
        existing = result.scalar_one_or_none()

        if existing:
            # Update fields
            update_data = event.model_dump(exclude_unset=True)
            for key, value in update_data.items():
                if key != "is_new":  # Don't overwrite is_new flag
                    setattr(existing, key, value)
            session.add(existing)
        else:
            session.add(event)
            count += 1

    await session.commit()
    return count
