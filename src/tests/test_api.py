import pytest
from datetime import datetime, timezone
from venvi.models.hackathon import Hackathon

@pytest.mark.asyncio
async def test_read_hackathons_empty(client):
    response = await client.get("/hackathons/")
    assert response.status_code == 200
    assert response.json() == []

@pytest.mark.asyncio
async def test_sync_endpoint_mocked(client, session):
    from unittest.mock import patch
    with patch("venvi.services.ingestion.fetch_euro_hackathons", return_value=[]):
        response = await client.post("/hackathons/sync")
        assert response.status_code == 200
        assert response.json() == {"message": "Sync complete", "new_items": 0}

@pytest.mark.asyncio
async def test_create_and_read_hackathon(client, session):
    hackathon = Hackathon(
        id="209e1cfb-5040-4086-bfec-f67bdc3380ff",
        name="API Test Hackathon",
        city="API City",
        country_code="AC",
        date_start=datetime(2026, 2, 3, 0, 0, 0, tzinfo=timezone.utc),
        date_end=datetime(2026, 2, 4, 0, 0, 0, tzinfo=timezone.utc),
        topics=[],
        url="https://api.test",
        status="upcoming"
    )
    session.add(hackathon)
    await session.commit()

    response = await client.get("/hackathons/")
    assert response.status_code == 200
    data = response.json()
    assert len(data) == 1
    assert data[0]["name"] == "API Test Hackathon"
