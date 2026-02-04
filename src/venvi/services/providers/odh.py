from datetime import datetime
from typing import Any

import httpx

from venvi.models.event import Event
from venvi.services.providers.base import BaseEventProvider


class ODHProvider(BaseEventProvider):
    @property
    def source_name(self) -> str:
        return "odh"

    async def fetch_events(self) -> list[dict[str, Any]]:
        api_url = "https://tourism.opendatahub.com/v1/Event"
        params = {"pagenumber": 1, "pagesize": 20}
        async with httpx.AsyncClient() as client:
            response = await client.get(api_url, params=params)
            response.raise_for_status()
            return response.json().get("Items", [])

    def map_event(self, raw: dict[str, Any]) -> Event:
        details = raw.get("Detail", {})

        def get_localized(obj: dict, key: str) -> str | None:
            for lang in ["en", "it", "de"]:
                val = obj.get(lang, {}).get(key)
                if val:
                    return val
            return None

        title = get_localized(details, "Title") or "Untitled Event"
        description = get_localized(details, "BaseText") or get_localized(
            details, "IntroText"
        )

        city = raw.get("ContactInfos", {}).get("en", {}).get("City") or "Unknown"

        image_url = None
        gallery = raw.get("ImageGallery", [])
        if gallery:
            image_url = gallery[0].get("ImageUrl")

        raw_id = raw.get("Id") or str(hash(title + str(raw.get("DateBegin"))))

        return Event(
            id=f"{self.source_name}:{raw_id}",
            title=title,
            description=description,
            date_start=datetime.fromisoformat(raw["DateBegin"])
            if raw.get("DateBegin")
            else datetime.now(),
            date_end=datetime.fromisoformat(raw["DateEnd"])
            if raw.get("DateEnd")
            else datetime.now(),
            location=city,
            url=f"https://opendatahub.com/events/{raw_id}",  # Placeholder URL
            image_url=image_url,
            source_name=self.source_name,
            source_id=raw_id,
            topics=[],
            category="general",
            is_new=True,
        )
