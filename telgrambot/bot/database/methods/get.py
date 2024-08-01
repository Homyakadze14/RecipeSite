from sqlalchemy import exc

from bot.database.main import Database
from bot.database.models.main import TgUsers


def get_tg_user_id(user_id):
    try:
        return Database().session.query(TgUsers).filter(TgUsers.user_id == user_id).one()
    except exc.NoResultFound:
        return None


def get_user_id(tg_user_id):
    try:
        return Database().session.query(TgUsers).filter(TgUsers.telegram_user_id == tg_user_id).one()
    except exc.NoResultFound:
        return None
