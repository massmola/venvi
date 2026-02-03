from datetime import datetime
from typing import Any

import httpx
from venvi.models.event import Event
from venvi.services.providers.base import BaseEventProvider


class EuroHackathonsProvider(BaseEventProvider):
    @property
    def source_name(self) -> str:
        return "euro_hackathons"

    async def fetch_events(self) -> list[dict[str, Any]]:
        api_url = "https://euro-hackathons.vercel.app/api/hackathons"
        async with httpx.AsyncClient() as client:
            response = await client.get(api_url, params={"status": "upcoming"})
            response.raise_for_status()
            return response.json().get("data", [])

    def map_event(self, raw: dict[str, Any]) -> Event:
        return Event(
            id=f"{self.source_name}:{raw['id']}",
            title=raw["name"],
            description=raw.get("notes"),
            date_start=datetime.fromisoformat(raw["date_start"]),
            date_end=datetime.fromisoformat(raw["date_end"]),
            location=f"{raw['city']}, {raw['country_code']}",
            url=raw["url"],
            image_url=None,
            source_name=self.source_name,
            source_id=raw["id"],
            topics=raw.get("topics", []),
            category="hackathon",
            is_new=True,
        )
