import routes

from fastapi import FastAPI
import uvicorn

from db import *

db = Python_db()
db.test_exec()

@asynccontextmanager
async def lifespan(app: FastAPI):
    await Python_db()
    yield

app = FastAPI(
    debug=True, # TODO: change
    title="chatAPI", 
)

for route in routes.__all__:
    router = getattr(routes, route)
    app.include_router(router)

