from fastapi import APIRouter
from .auth import router as auth_router
from .api import router as axios_router
from .master import router as master_router
from .setting import router as setting_router
from .config import router as config_router
from .Test import router as test_router

api_router = APIRouter()

api_router.include_router(auth_router, prefix="/auth", tags=["Auth"])
api_router.include_router(axios_router, prefix="/api", tags=["Api"])
api_router.include_router(master_router, prefix="/master", tags=["Master"])
api_router.include_router(setting_router, prefix="/setting", tags=["Setting"])
api_router.include_router(config_router, prefix="/config", tags=["Config"])
api_router.include_router(test_router, prefix="/test", tags=["Test"])