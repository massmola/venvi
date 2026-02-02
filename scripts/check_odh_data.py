import asyncio

from sqlalchemy.ext.asyncio import async_sessionmaker
from sqlmodel import select
from venvi.core.db import engine
from venvi.models.odh_event import ODHEvent


async def check_data() -> None:
    # Use class_ keyword properly or use correct sessionmaker
    async_session = async_sessionmaker(engine, expire_on_commit=False)
    async with async_session() as session:
        result = await session.execute(select(ODHEvent))
        events = result.scalars().all()
        print(f"Found {len(events)} ODH events in the database.")
        for event in events[:5]:
            print(f"- {event.title} ({event.date_start})")


if __name__ == "__main__":
    asyncio.run(check_data())
