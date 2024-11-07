from fastapi import APIRouter

router = APIRouter(
    prefix="/ping", 
    tags=["Ping"],
)

@router.get("/", status_code=200)
def ping():
    """
    Ping to test the connection.
    """
    return "pong"