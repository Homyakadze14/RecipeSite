from bot.handlers import other, user
from aiogram import Router


def get_all_routers() -> list[Router]:
    return [user.main.router, other.router]
