import asyncio

import httpx


async def test_taken_status() -> None:
    async with httpx.AsyncClient(base_url="http://127.0.0.1:8000") as client:
        # 1. Get an ODH event
        print("Fetching ODH events...")
        response = await client.get("/events/odh")
        if response.status_code != 200:
            print("Failed to fetch events")
            return

        events = response.json()
        if not events:
            print("No events found to test with.")
            # Trigger sync
            print("Triggering sync...")
            await client.post("/events/odh/sync")
            events = (await client.get("/events/odh")).json()

        if not events:
            print("Still no events.")
            return

        event_id = events[0]["id"]
        print(f"Testing with event ID: {event_id}")

        # 2. Mark as taken
        print("Marking as taken...")
        response = await client.patch(
            f"/events/odh/{event_id}/taken", params={"taken": "true"}
        )
        print(f"Status Code: {response.status_code}")
        print(f"Response: {response.json()}")
        assert response.json()["taken"]

        # 3. Verify filter
        print("Verifying filter (taken=true)...")
        response = await client.get("/events/odh", params={"taken": "true"})
        filtered = response.json()
        assert any(e["id"] == event_id for e in filtered)
        print("Filter 'taken=true' verified.")

        # 4. Mark as not taken
        print("Marking as NOT taken...")
        response = await client.patch(
            f"/events/odh/{event_id}/taken", params={"taken": "false"}
        )
        assert not response.json()["taken"]
        print("Marked as not taken.")


if __name__ == "__main__":
    try:
        asyncio.run(test_taken_status())
    except Exception as e:
        print(f"Error: {e}")
