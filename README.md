# go-rti-testing
Полное описание тестового задания: [go-rti-testing](https://github.com/vafinvr/go-rti-testing)  
Приложение тестировалось на _Ubuntu 18.04.4 LTS_

## Installation
1. Убедиться, что уже установлен [Docker Engine](https://docs.docker.com/get-docker/);
2. Собрать _Docker_ образ: _docker build -t go-rti-testing ._;
3. Создать и запустить контейнер: _docker run --publish 8000:8080 --detach --name grt go-rti-testing_;
4. Убедившись в работе, можно удалить контейнер: _docker rm --force grt_.

## Other

Пример запроса

```json
curl --location --request POST 'http://localhost:8000/calculate' \
--header 'Content-Type: application/json' \
--data-raw '{
   "product":{
      "name":"Игровой",
      "components":[
         {
            "isMain":true,
            "name":"Интернет",
            "prices":[
               {
                  "cost":500,
                  "priceType":"COST",
                  "ruleApplicabilities":[
                     {
                        "codeName":"technology",
                        "operator":"EQ",
                        "value":"xpon"
                     },
                     {
                        "codeName":"internetSpeed",
                        "operator":"EQ",
                        "value":"100"
                     }
                  ]
               },
               {
                  "cost":900,
                  "priceType":"COST",
                  "ruleApplicabilities":[
                     {
                        "codeName":"technology",
                        "operator":"EQ",
                        "value":"xpon"
                     },
                     {
                        "codeName":"internetSpeed",
                        "operator":"EQ",
                        "value":"200"
                     }
                  ]
               },
               {
                  "cost":10,
                  "priceType":"DISCOUNT",
                  "ruleApplicabilities":[
                     {
                        "codeName":"internetSpeed",
                        "operator":"GTE",
                        "value":"50"
                     }
                  ]
               }
            ]
         },
         {
            "name":"ADSL Модем",
            "prices":[
               {
                  "cost":300,
                  "priceType":"COST",
                  "ruleApplicabilities":[
                     {
                        "codeName":"technology",
                        "operator":"EQ",
                        "value":"adsl"
                     }
                  ]
               }
            ]
         }
      ]
   },
   "conditions":[
      {
         "ruleName":"technology",
         "value":"xpon"
      },
      {
         "ruleName":"internetSpeed",
         "value":"200"
      }
   ]
}'
```