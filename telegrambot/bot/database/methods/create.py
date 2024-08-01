from sqlalchemy import exc

from bot.database.main import Database
from loguru import logger


def create_tg_user(tg_user):
    session = Database().session
    session.add(tg_user)

    try:
        session.commit()
    except exc.IntegrityError:
        logger.error(f"I can't create user with tg id {tg_user.telegram_user_id}!")
        return session.rollback()
