from aiogram import Router, F
from aiogram.filters import Command
from aiogram.fsm.context import FSMContext
from aiogram.types import Message
from aiogram.fsm.state import StatesGroup, State

from bot.database.methods import other, create, delete
from bot.database.models.main import TgUsers

from os import environ
import requests
from loguru import logger

router = Router()


@router.message(Command('start'))
async def start(message: Message, state: FSMContext):
    await state.clear()
    if other.user_exist(message.from_user.id):
        await message.answer(f"Привет, <b>{message.from_user.first_name}</b>!\n"
                             f"Как только пользователь, на которого вы подписаны, выложит рецепт, я вам сообщу")
    else:
        await message.answer(f"Привет, <b>{message.from_user.first_name}</b>!\n"
                             f"Пожалуйста, <a href="">сгенерируйте специальный код</a> для входа на сайте и пришлите его мне, "
                             f"использовав команду /login")


class Login(StatesGroup):
    token = State()


@router.message(Command('login'))
async def login(message: Message, state: FSMContext):
    await state.clear()
    if other.user_exist(message.from_user.id):
        await message.answer(f"Вы уже вошли под своим аккаунтом, если хотите выйти напишите /logout")
    else:
        await message.answer(f"Пришлите мне токен, который вы сгенерировали: ")
        await state.set_state(Login.token)


@router.message(F.text, Login.token)
async def enter_token(message: Message, state: FSMContext):
    try:
        r = requests.post(environ.get("BACKEND_BASE_URL") + "/auth/checktgtoken", json={"token": message.text})
        if r.status_code == 400:
            await message.answer(f"Ваш код не действителен, попробуйте ещё раз")
            return
        elif r.status_code == 200:
            tgUser = TgUsers()
            tgUser.user_id = r.json()['user_id']
            tgUser.telegram_user_id = message.from_user.id
            create.create_tg_user(tgUser)
            await message.answer(f"Вы успешно вошли в свой аккаунт!")
            await state.clear()
        else:
            await message.answer(
                f"Упс! Возникла ошибка на сервере, попробуйте снова через небольшой промежуток времени!")
            await state.clear()
    except Exception as e:
        logger.error(e)
        await message.answer(f"Упс! Возникла ошибка на сервере, попробуйте снова через небольшой промежуток времени!")
        await state.clear()


@router.message(Command('logout'))
async def logout(message: Message, state: FSMContext):
    await state.clear()
    if other.user_exist(message.from_user.id):
        tgUser = TgUsers()
        tgUser.telegram_user_id = message.from_user.id
        delete.delete_tg_user(tgUser)
        await message.answer(f"Вы успешно вышли из своего аккаунта!\nЕсли хотите войти, напишите /login")
    else:
        await message.answer(f"Вы не вошли в аккаунт")
