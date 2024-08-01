from bot.database.methods import get


def user_exist(tg_user_id):
    return get.get_user_id(tg_user_id) is not None
