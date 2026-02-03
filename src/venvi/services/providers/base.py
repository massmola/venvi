from abc import ABC, abstractmethod
from typing import Any

from sqlalchemy.ext.asyncio import AsyncSession
from venvi.models.event import Event


class BaseEventProvider(ABC):
    """
    Abstract base class for event providers.
    """

    @property
    @abstractmethod
    def source_name(self) -> str:
        """The identifier for this source (e.g., 'euro_hackathons')."""
        pass

    @abstractmethod
    async def fetch_events(self) -> list[dict[str, Any]]:
        """Fetch raw event data from the source."""
        pass

    @abstractmethod
    def map_event(self, raw: dict[str, Any]) -> Event:
        """Map raw data to the unified Event model."""
        pass

    async def sync_events(self, session: AsyncSession) -> int:
        """Fetch, map, and synchronize events with the database."""
        from sqlmodel import select

        raw_events = await self.fetch_events()
        count = 0

        for raw in raw_events:
            event = self.map_event(raw)
            # Ensure source info is set
            event.source_name = self.source_name

            # Check if exists
            statement = select(Event).where(Event.id == event.id)
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
