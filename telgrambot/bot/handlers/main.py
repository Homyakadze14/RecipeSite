from bot.handlers import admin, other, user
from aiogram import Router


def get_all_routers() -> list[Router]:
    return [admin.main.router, user.main.router, other.router]
