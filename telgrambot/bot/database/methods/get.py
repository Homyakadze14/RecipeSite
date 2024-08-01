from sqlalchemy import exc

from bot.database.main import Database
from bot.database.models.main import TgUsers, Subscriptions


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


def get_subscribers(creator_id):
    try:
        return Database().session.query(Subscriptions).filter(Subscriptions.creator_id == creator_id).all()
    except exc.NoResultFound:
        return None