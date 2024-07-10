# RecipeSite

## Регистрация
URL: api/v1/signup
METHOD: POST
Принимает форму с полями:
icon - фотография,
email - текст
login - текст(мин 3 симв, макс 20 симв),
password - текст(мин 8 симв, макс 50 симв),
about - текст(макс 1500 симв)

Коды:
200,
400,
500

Ставит куку

## Логин
URL: api/v1/signin
METHOD: POST
Принимает json с полями:
email - текст
ИЛИ
login - текст,
password - текст

Коды:
200,
400,
500

Ставит куку

## Логоут
URL: api/v1/logout
METHOD: POST

Коды:
200,
500

Ставит куку