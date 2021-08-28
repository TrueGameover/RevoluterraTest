1. `cp .env.dist .env`
2. `docker-compose up --build`
3. http://localhost:5000/swagger/index.html

В контейнере поймал проблему связанную с https когда mtu явно не задан для докер сервиса в systemd.
В этом кейсе при подключении к яндексу падает в "tls handshake timeout".
