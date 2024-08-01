from aiogram import Bot, Dispatcher
from aiogram.client.default import DefaultBotProperties
from aiogram.fsm.storage.memory import MemoryStorage

from loguru import logger

from bot.database.models import register_models

from bot.handlers.main import get_all_routers

from dotenv import load_dotenv
import os


async def start_bot():
    if os.environ.get("GIN_MODE") != "release":
       load_dotenv(os.path.join(os.path.dirname(__file__), '.env'))

    bot = Bot(token=os.environ.get("TOKEN"), default=DefaultBotProperties(parse_mode='HTML'))
    dp = Dispatcher(storage=MemoryStorage())

    logger.info('Bot starts!')
    dp.include_routers(*get_all_routers())
    register_models()

    await bot.delete_webhook(drop_pending_updates=True)
    await dp.start_polling(bot)
