from sqlalchemy import exc

from bot.database.main import Database
from bot.database.models.main import TgUsers
from loguru import logger


def delete_tg_user(tg_user):
    session = Database().session
    session.query(TgUsers).filter(TgUsers.telegram_user_id == tg_user.telegram_user_id).delete()
    try:
        session.commit()
    except exc.IntegrityError:
        logger.error(f"I can't delete user with tg id {tg_user.telegram_user_id}!")
        return session.rollback()
