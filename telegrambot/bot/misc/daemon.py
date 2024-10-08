import aio_pika, asyncio, os
from loguru import logger
from os import environ
import ast, requests
from bot.database.methods import get
from aiogram.utils.markdown import link

def get_post_info(id):
    try:
        url = environ.get("BACKEND_BASE_URL") + f"/recipe/{id}"
        recipe_url = environ.get("RECIPE_URL") + f"{id}"
        r = requests.get(url)
        if r.status_code == 200:
            recipe = r.json()['info']['recipe']
            author = r.json()['info']['recipe']['author']
            info = (f"*Новый рецепт\\!*\n\n*Название:* {recipe['title']}\n*Описание:* {recipe['about']}\n"
                    f"*Создатель:* {author['login']}\n" +
                    f"_{link('Подробнее', recipe_url)}_")
            return info
        else:
            logger.error("Server error")
    except Exception as e:
        logger.error(e)
        return ""


async def send_messages(bot, message):
    for subscriber in get.get_subscribers(message['CreatorID']):
        tg_user = get.get_tg_user_id(subscriber.subscriber_id)
        if tg_user is not None:
            try:
                post = get_post_info(message['RecipeID'])
                if post == "":
                    continue
                await bot.send_message(chat_id=tg_user.telegram_user_id, text=post, parse_mode="MarkdownV2")
            except Exception as e:
                logger.error(e)


async def run(bot, loop):
    counter = 10
    for i in range(1, 11):
        try:
            connection = await aio_pika.connect_robust(
                environ.get("RMQ_URL"), loop=loop
            )
        except:
            logger.error(f"Try to connect to rabbit: {counter}")
            counter -= 1
            await asyncio.sleep(5)

    if counter == 0:
        logger.error(f"Can't connect to rabbit")
        os._exit(1)

    logger.info(f"Connect to rabbit")

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
                    
                    logger.info(f"Get message")

                    await send_messages(bot, msg)