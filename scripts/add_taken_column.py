import asyncio

from sqlalchemy import text
from venvi.core.db import engine


async def migrate() -> None:
    async with engine.begin() as conn:
        print("Migrating Hackathon table...")
        try:
            await conn.execute(
                text("ALTER TABLE hackathon ADD COLUMN taken BOOLEAN DEFAULT FALSE")
            )
            print("Hackathon table updated.")
        except Exception as e:
            print(f"Hackathon table update failed (maybe column exists?): {e}")

        print("Migrating ODHEvent table...")
        try:
            await conn.execute(
                text("ALTER TABLE odhevent ADD COLUMN taken BOOLEAN DEFAULT FALSE")
            )
            print("ODHEvent table updated.")
        except Exception as e:
            print(f"ODHEvent table update failed (maybe column exists?): {e}")


if __name__ == "__main__":
    asyncio.run(migrate())
