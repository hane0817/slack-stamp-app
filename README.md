#### アプリケーション名

Slack-stamp-app

#### アプリーションの機能

Slackのリアクション機能で用いるための文字画像を生成するアプリケーション

#### 対応言語

- 日本語
- 英語
- 中国語

#### 環境構築

- ファイル操作のエラーが発生しないような場所(PATH)でclone
- cd frontend
- npm install
- cd ../
- docker-compose up

#### 注意点

 - 初めてdocker-compose upを行った場合、backendのサーバとdbの起動がずれることで接続失敗する可能性があるので、一度docker-compose stop後に再度docker-compose upをすると解決する