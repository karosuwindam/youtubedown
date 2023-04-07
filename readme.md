# Youtubeのダウンロード補助アプリ

## 概要
Web待機を実施して、対象サーバ

pythonのツールであるyoutube-dlを利用してmp3に変換する


## 主なAPIについて

|url|機能|備考|
|--|--|--|
|/view/:id|対象のmp3ファイルのID3v2タグを表示||
|/edit/:id|対象のmp3ファイルのID3v2タグを編集|未完成|
|/mp3/:id|対象のmp3ファイルをダウンロード||
|/mp3image/:id|対象のmp3ファイルのタグから画像表示||
|/download|YouTubeのURLを登録することでmp3に変換する||
|/list|ダウンロードできるmp3のファイル名リスト||
|/health|YouTubeのダウンロード状態を確認||