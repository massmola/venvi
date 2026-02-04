from contextlib import asynccontextmanager
from pathlib import Path

from fastapi import FastAPI
from fastapi.staticfiles import StaticFiles

import venvi.api.routers.events
from venvi.core.db import init_db
from venvi.web.router import router as web_router


@asynccontextmanager
async def lifespan(app: FastAPI):
    await init_db()
    yield


app = FastAPI(
    title="Venvi - EU Event Suggestion Platform",
    description="Discover and sync events from multiple sources",
    version="0.1.0",
    lifespan=lifespan,
)

# Mount Static Files
BASE_DIR = Path(__file__).resolve().parent
STATIC_DIR = BASE_DIR / "static"
STATIC_DIR.mkdir(parents=True, exist_ok=True)
app.mount("/static", StaticFiles(directory=str(STATIC_DIR)), name="static")

# Include Routers
app.include_router(web_router)
app.include_router(venvi.api.routers.events.router)
