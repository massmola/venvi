from sqlalchemy.ext.asyncio import AsyncSession

from venvi.services.providers.euro_hackathons import EuroHackathonsProvider
from venvi.services.providers.odh import ODHProvider

PROVIDERS = [
    EuroHackathonsProvider(),
    ODHProvider(),
]


async def sync_all_events(session: AsyncSession) -> dict[str, int]:
    """
    Synchronizes events from all registered providers.

    Args:
        session: The asynchronous database session.

    Returns:
        dict: A mapping of provider source names to the count of new events added.
    """
    results = {}
    for provider in PROVIDERS:
        count = await provider.sync_events(session)
        results[provider.source_name] = count
    return results
