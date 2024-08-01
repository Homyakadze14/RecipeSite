from sqlalchemy import Column, Integer

from bot.database.main import Database


class TgUsers(Database.BASE):
    __tablename__ = 'tgusers'
    id = Column(Integer, primary_key=True)
    user_id = Column(Integer, nullable=False, unique=True)
    telegram_user_id = Column(Integer, nullable=False, unique=True)


class Subscriptions(Database.BASE):
    __tablename__ = 'subscriptions'
    id = Column(Integer, primary_key=True)
    creator_id = Column(Integer, nullable=False)
    subscriber_id = Column(Integer, nullable=False)


def register_models():
    Database.BASE.metadata.create_all(Database().engine)
