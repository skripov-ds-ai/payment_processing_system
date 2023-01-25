# Микросервис для работы с балансом пользователей
![Project language][badge_language]
[![codecov](https://codecov.io/gh/nizhikebinesi/payment_processing_system/branch/main/graph/badge.svg?token=ZLTO2VGPPH)](https://codecov.io/gh/nizhikebinesi/payment_processing_system)

[badge_language]:https://img.shields.io/badge/language-go_1.19-blue.svg?longCache=true

## Цели
1. Сделать без генератора из swagger-документации
2. Унифицировать способ взаимодействия с БД, чтобы была возможность 
подключить MySQL(sql из stdlib) и postgres(pgx, не стандартный интерфейс)
3. Добавить Swagger-документацию
4. Добавить генератор C4-диаграмм