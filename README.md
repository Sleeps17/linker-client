# Консольный Linux клиент для сервиса linker
Сервис linker - ``https://github.com/Sleeps17/linker``.

## Установка
1)   ``git@github.com:Sleeps17/linker-client.git``
2)   ``sudo snap install task --classic``
3)   ``task install``

## Документация по использованию
### Работа с топиками
1) Создание топика: ``linker post_topic --topic=some_topic``
2) Удаление топика: ``linker delete_topic --topic=some_topic``
3) Просмотр спика топиков: ``linker list_topics``
### Работа со ссылками
1) Добавление ссылки в топик: ``linker post_link --topic=some_topic --link=some_link --alias=some_alias``
2) Получение ссылки по ее алиасу: ``linker pick_link --topic=some_topic --alias=some_alias``
3) Удаление ссылки по ее алиасу: ``linker delete_link --topic=some_topic --alias=some_alias``
4) Получение списка всех ссылок расположенных в топике: ``linker list_links --topic=some_topic``

### Если вам наравится суть сервиса и вы хотите его использовать, то можете самомтоятельно написать клиент:
protobuf files: ``https://github.com/Sleeps17/linker-protos``
