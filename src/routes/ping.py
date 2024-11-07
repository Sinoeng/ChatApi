from fastapi import APIRouter

router = APIRouter(
    prefix="/ping", 
    tags=["Ping"],
)

@router.get("/", status_code=200)
def ping() -> str:
    """
    Ping to test the connection.
    """
    return "pong"

@router.get("/db", status_code=201)
def addVar():
    return db.exec("SHOW DATABASES;")

