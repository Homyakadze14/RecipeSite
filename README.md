# RecipeSite

## Регистрация
**URL:** api/v1/signup<br>
**METHOD:** POST<br>
**Принимает json с полями:**<br>
email - текст<br>
login - текст(мин 3 симв, макс 20 симв)<br>
password - текст(мин 8 симв, макс 50 симв)<br>
<br>
**Коды:**<br>
200,<br>
400,<br>
500<br>
<br>
**Ставит куку**

## Логин
**URL:** api/v1/signin<br>
**METHOD:** POST<br>
**Принимает json с полями:**<br>
email - текст<br>
ИЛИ<br>
login - текст<br>
password - текст<br>
<br>
**Коды:**<br>
200,<br>
400,<br>
500<br>
<br>
**Ставит куку**

## Логаут
**URL:** api/v1/logout<br>
**METHOD:** POST<br>
<br>
**Коды:**<br>
200,<br>
401,<br>
500<br>
<br>
**Ставит куку**


## Обновление login, email, about, icon
**URL:** api/v1/user/{login}<br>
**METHOD:** PUT<br>
**Принимает форму с полями:**<br>
icon - текст<br>
email - текст<br>
login - текст(мин 3 симв, макс 20 симв)<br>
about - текст(макс 1500  симв)<br>
<br>
**Коды:**<br>
200,<br>
401,<br>
404,<br>
500<br>

## Обновление password
**URL:** api/v1/user/{login}/password<br>
**METHOD:** PUT<br>
**Принимает json с полями:**<br>
password - текст(мин 8 симв, макс 50 симв)<br>
<br>
**Коды:**<br>
200,<br>
401,<br>
404,<br>
500<br>
**Удаляет все активные сессии**

## Получение всех рецептов
**URL:** api/v1/recipe<br>
**METHOD:** GET<br>
**Коды:**<br>
200,<br>
500<br>
**Отдаёт все рецепты в json**

## Получение всех рецептов с фильтрацией
**URL:** api/v1/recipe<br>
**METHOD:** POST<br>
**Принимает json с полями:**<br>
limit - число (количество рецептов из бд)<br>
offset - число (сколько элементов пропустить)<br>
query - текст (ищет все возможные совпадения в title, about, ingridients)<br>
order_field - текст (сортирует по title, complexitiy, updated_at)<br>
order_by - число от -1 до 1 (-1 - по возрастанию, 0 - никак, 1 - по убыванию)<br>
P.S Все эти поля опциональны
<br>
**Коды:**<br>
200,<br>
500<br>
**Отдаёт отфильтрованные рецепты в json**

## Получение рецепта
**URL:** api/v1/recipe/{id:[0-9]+}<br>
**METHOD:** GET<br>
<br>
**Коды:**<br>
200,<br>
404,<br>
500<br>
**Отдаёт рецепт в json**

## Создание рецепта
**URL:** api/v1/user/{login}/recipe<br>
**METHOD:** POST<br>
**Принимает форму с полями:**<br>
title - текст(мин 3 симв, макс 50 симв)<br>
about - текст(макс 2500 симв)<br>
complexitiy - число от 1 до 3 (Сложность приготовления лёгкая\средняя\тяжёлая)<br>
need_time - текст (сколько времени нужно для приготовления)<br>
ingridients - текст (макс 1500 симв) (какие ингридиенты нужны)<br>
photos - изображения<br>
<br>
**Коды:**<br>
200,<br>
401,<br>
404,<br>
500<br>

## Обновление рецепта
**URL:** api/v1/user/{login}/recipe/{id:[0-9]+}<br>
**METHOD:** PUT<br>
**Принимает форму с полями:**<br>
title - текст(мин 3 симв, макс 50 симв)<br>
about - текст(макс 2500 симв)<br>
complexitiy - число от 1 до 3 (Сложность приготовления лёгкая\средняя\тяжёлая)<br>
need_time - текст (сколько времени нужно для приготовления)<br>
ingridients - текст (макс 1500 симв) (какие ингридиенты нужны)<br>
photos - изображения<br>
<br>
**Коды:**<br>
200,<br>
401,<br>
404,<br>
500<br>

## Удаление рецепта
**URL:** api/v1/user/{login}/recipe/{id:[0-9]+}<br>
**METHOD:** DELETE<br>
<br>
**Коды:**<br>
200,<br>
401,<br>
404,<br>
500<br>

# Лайкнуть рецепт
**URL:** api/v1/recipe/{id:[0-9]+}/like<br>
**METHOD:** POST<br>
<br>
**Коды:**<br>
200,<br>
401,<br>
500<br>

# Убрать лайк
**URL:** api/v1/recipe/{id:[0-9]+}/unlike<br>
**METHOD:** POST<br>
<br>
**Коды:**<br>
200,<br>
401,<br>
500<br>