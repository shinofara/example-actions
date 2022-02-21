# コードを実行するコンテナイメージ
FROM golang:1.7

COPY . .
RUN go build -o /bin/check
ENTRYPOINT ["/entrypoint.sh"]