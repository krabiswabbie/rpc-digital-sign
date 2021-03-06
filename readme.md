# Пример RPC-сервера с использованием цифровой подписи

## Задание
Написать на golang, ping-pong сервер, клиент. 

Клиент (консольная утилита) должен отправлять пинг запрос на сервер, подписывая запрос приватным ключом и проверять подпись в ответе. Сервер должен обрабатывать входящий `ping` запрос, проверять подпись созданную на клиенте и отправить подписанный ответ. 

- В качестве зависимости использовать `github.com/ethereum/go-ethereum v1.9.7`
- Для взаимодействия сервера и клиента использовать RPC `github.com/ethereum/go-ethereum/rpc` 
- Для вызова из терминала использовать `gopkg.in/urfave/cli.v1` 
- Для работы с подписями использовать `github.com/ethereum/go-ethereum/crypto` или криптопро (криптопро предпочтительнее)

## Запуск

```bash
    ## Server
    ./ether server

    ## Client
    ./ether client
```

## Пояснения к реализации

Текущая реализация показывает базовый функционал JSON-RPC (`go-ethereum/rpc` на `unix sockets`) и цифровой подписи (`go-ethereum/crypto`). Код может быть улучшен следующим образом:

1. Обмен публичными ключами должен выполняться один раз при подключении клиента (сейчас ключ передается в каждом сообщении).
2. Предусмотреть корректное закрытие сервера (и отключение клиента).

## Traps and pitfalls

### RPC-методы

Вызов нужного RPC-метода осуществляется следующим образом:

```
func callEcho(c *rpc.Client, msg *EchoMsg) (*EchoMsg, error) {
	resp := new(EchoMsg)
	err := c.Call(resp, "service_echo", msg)
	return resp, err
}
```

Здесь `service_echo` - составное (полное) имя метода. `service` - имя модуля, `echo` - имя метода. Имя метода в данной записи должно быть с маленькой буквы(!), в то время как его объявление в коде - с большой.

```
func (s *TestService) Echo(msg *EchoMsg) (*EchoMsg, error) {
    ...
}
```


### Длина (формат) подписи
```
func Decode(msg *EchoMsg) bool {
	sign := msg.Signature[:len(msg.Signature)-1] // remove recovery id
	hash := eth.Keccak256(msg.Data)
	return eth.VerifySignature(msg.PubKey, hash, sign)
}
```
Функция `VerifySignature` принимает подпись из 64 байт (`The signature should have the 64 byte [R || S] format.`), в то время как полная длина подписи, возращаемой функцией `Sign()` 65 байт (`The produced signature is in the [R || S || V] format where V is 0 or 1.`). Лишним байтом (`V`) можно пренебречь)

