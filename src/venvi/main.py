from contextlib import asynccontextmanager
from pathlib import Path

from fastapi import FastAPI
from fastapi.staticfiles import StaticFiles

from venvi.api.routers import hackathons
from venvi.core.db import init_db
from venvi.web.router import router as web_router


@asynccontextmanager
async def lifespan(app: FastAPI):
    await init_db()
    yield


app = FastAPI(
    title="Venvi - Euro Hackathons",
    description="Weekend trip suggestion turned Hackathon Aggregator",
    version="0.1.0",
    lifespan=lifespan,
)

# Mount Static Files
BASE_DIR = Path(__file__).resolve().parent
app.mount("/static", StaticFiles(directory=str(BASE_DIR / "static")), name="static")

# Include Routers
app.include_router(web_router)
app.include_router(hackathons.router)
