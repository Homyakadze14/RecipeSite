from sqlalchemy import Column, Integer

from bot.database.main import Database


class TgUsers(Database.BASE):
    __tablename__ = 'tgusers'
    id = Column(Integer, primary_key=True)
    user_id = Column(Integer, nullable=False, unique=True)
    telegram_user_id = Column(Integer, nullable=False, unique=True)


def register_models():
    Database.BASE.metadata.create_all(Database().engine)
