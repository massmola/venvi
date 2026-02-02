from contextlib import asynccontextmanager

from fastapi import FastAPI

from venvi.api.routers import hackathons
from venvi.core.db import init_db


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

app.include_router(hackathons.router)


@app.get("/")
async def root() -> dict[str, str]:
    return {"message": "Welcome to Venvi! Check /docs for API."}
