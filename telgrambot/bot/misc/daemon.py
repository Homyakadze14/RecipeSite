import aio_pika
from os import environ
import ast, requests
from bot.database.methods import get
from aiogram.utils.markdown import link

def get_post_info(id):
    try:
        url = environ.get("BACKEND_BASE_URL") + f"/recipe/{id}"
        r = requests.get(url)
        if r.status_code == 200:
            recipe = r.json()['info']['recipe']
            author = r.json()['info']['author']
            info = (f"*Новый рецепт\\!*\n\n*Название:* {recipe['title']}\n*Описание:* {recipe['about']}\n"
                    f"*Создатель:* {author['login']}\n" +
                    f"_{link('Подробнее', url)}_")
            return info
        else:
            print("Server error")
    except Exception as e:
        print(e)
        return ""


async def send_messages(bot, message):
    for subscribers in get.get_subscribers(message['CreatorID']):
        tg_user = get.get_tg_user_id(subscribers.subscriber_id)
        if tg_user is not None:
            try:
                post = get_post_info(message['RecipeID'])
                if post == "":
                    continue
                await bot.send_message(chat_id=tg_user.telegram_user_id, text=post, parse_mode="MarkdownV2")
            except Exception as e:
                print(e)


async def run(bot, loop):
    connection = await aio_pika.connect_robust(
        environ.get("RMQ_URL"), loop=loop
    )

    async with connection:
        queue_name = "new_recipe"

        channel: aio_pika.abc.AbstractChannel = await connection.channel()

        queue: aio_pika.abc.AbstractQueue = await channel.declare_queue(
            queue_name,
        )

        async with queue.iterator() as queue_iter:
            async for message in queue_iter:
                async with message.process():
                    msg = ast.literal_eval(message.body.decode('utf-8'))

                    if queue.name in message.body.decode():
                        break

                    await send_messages(bot, msg)