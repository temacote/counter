Сервер с функцией сохранения факта обработки запросов в память 

собрать бинарь можно выполнив make в корне 

Сервер состоит из двух запускаемых независимо инстансов - grpc сервера и http gateway сервера,
который прокидывает запросы на grpc эндпоинты 

Список GRPC эндпоинтов объявляется в proto файле proto/common.proto
Генератор в generate.go
Имплементация в cmd/grpc/listeners/public

Основной конфиг сервиса лежит в bin/config.dev.yaml

В работе сервера имплементирована логика взаимодействия с:
 1. consul - для сервис дискавери и healthcheck
 2. jaegertracer - для трассировки запросов
 3. redis - для реализации IMDB схемы хранения данных по запросам
 4. pushgateway/prometheus - для сбора метрик

Для функционирования сервера необходимо (запуск grpc и http производить в разных терменалах):
 1. Собрать окружение docker-compose up -d в корне
 2. Запустить grpc сервер: go run sber_cloud/tw grpc_public -c /bin/config.dev.yaml --swagger proto/grpc_public.swagger.json
 3. Запустить http gateway сервер: go run sber_cloud/tw http -c /bin/config.dev.yaml
 
Теперь после окончания обработки запроса GET localhost:3018/v1/count (имплементирован в cmd/grpc/listeners/public/count_v1.go)
Информация по запросу будет сохранена в redis в обработчике формирования метадаты ответа cmd/http/gateway/response_handler.go
Для каждой записи в redis выставляется ttl в 60 секунд

При получении сервером сигнала SIGTERM происходит сохранение данных из редиса в файл cmd/http_server.go:159

При перезапуске сервиса происходит загрузка данных с файла на диске обратно в редис cmd/definition/http/gateway/response_handler.go:30

Окружение:
    1. consul - http://127.0.0.1:8500/ui/dc1/services
    2. jaeger - http://127.0.0.1:16686/search
    3. swagger ui - http://localhost:3018/swagger-ui/