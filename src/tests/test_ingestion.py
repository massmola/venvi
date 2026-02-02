from unittest.mock import patch, MagicMock
import pytest
from venvi.services.ingestion import sync_hackathons
from venvi.models.hackathon import Hackathon
from sqlmodel import select

SAMPLE_DATA = [
    {
        "id": "209e1cfb-5040-4086-bfec-f67bdc3380ff",
        "name": "Test Hackathon",
        "city": "Test City",
        "country_code": "TC",
        "date_start": "2026-02-03T00:00:00+00:00",
        "date_end": "2026-02-04T00:00:00+00:00",
        "topics": [],
        "notes": None,
        "url": "https://example.com",
        "status": "upcoming",
        "is_new": False
    }
]

@pytest.mark.asyncio
async def test_sync_hackathons(session):
    with patch("venvi.services.ingestion.fetch_euro_hackathons", return_value=SAMPLE_DATA) as mock_fetch:
        count = await sync_hackathons(session)
        assert count == 1
        
        result = await session.execute(select(Hackathon))
        hackathons = result.scalars().all()
        assert len(hackathons) == 1
        assert hackathons[0].name == "Test Hackathon"

@pytest.mark.asyncio
async def test_sync_hackathons_update(session):
    # Initial sync
    with patch("venvi.services.ingestion.fetch_euro_hackathons", return_value=SAMPLE_DATA):
        await sync_hackathons(session)
    
    # Update data
    updated_data = list(SAMPLE_DATA)
    updated_data[0]["name"] = "Updated Name"
    
    with patch("venvi.services.ingestion.fetch_euro_hackathons", return_value=updated_data):
        count = await sync_hackathons(session)
        assert count == 0 # No new items
        
        result = await session.execute(select(Hackathon))
        hackathon = result.scalars().first()
        assert hackathon.name == "Updated Name"
