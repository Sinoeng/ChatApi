import routes

from fastapi import FastAPI
import uvicorn

app = FastAPI(
    debug=True, 
    title="chatAPI", 
)

for route in routes.__all__:
    router = getattr(routes, route)
    app.include_router(router)