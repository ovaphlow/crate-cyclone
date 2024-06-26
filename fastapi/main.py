import os

from dotenv import load_dotenv
from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse
from starlette.exceptions import HTTPException as StarletteHTTPException

from application.schema_routes import router as schema_router
from application.setting_router import router as setting_router
from application.subscriber_routes import router as subscriber_router

load_dotenv()

app = FastAPI()


@app.exception_handler(StarletteHTTPException)
async def http_exception_handler(request: Request, exc: StarletteHTTPException):
    if exc.status_code == 500:
        title = "服务器错误"
    else:
        title = "异常错误"

    return JSONResponse(
        status_code=exc.status_code,
        content=dict(
            type="about:blank",
            status=exc.status_code,
            title=title,
            detail=str(exc.detail),
            instance=str(request.url)
        )
    )


@app.get("/")
async def root():
    return {"message": "Hello World"}


app.include_router(setting_router, prefix="/crate-api/setting", tags=["setting"])

app.include_router(subscriber_router, prefix="/crate-api/subscriber", tags=["subscriber"])

app.include_router(schema_router, prefix="/crate-api/database", tags=["schema"])

if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app, host="0.0.0.0", port=os.getenv("PORT", 8421))
