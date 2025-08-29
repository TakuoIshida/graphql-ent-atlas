# ビルドステージ
FROM golang:1.21-alpine AS builder

# 作業ディレクトリの設定
WORKDIR /app

# Go modulesファイルをコピー
COPY go.mod go.sum ./

# 依存関係のダウンロード
RUN go mod download

# アプリケーションソースコードをコピー
COPY . .

# アプリケーションのビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# 実行ステージ
FROM alpine:latest

# 証明書とタイムゾーン情報を追加
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# ビルドしたバイナリをコピー
COPY --from=builder /app/server .

# ポートの公開
EXPOSE 8080

# アプリケーションの実行
CMD ["./server"]