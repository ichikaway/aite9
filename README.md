# aite9 (port close check scanner)

サーバのポートが空いてないことを確認するポートスキャナー  
ポートが空いているのを検知するとエラー出力して終了ステータスは1以上で終わります。  
オプションでエラーが起こった場合にslack通知します。  

Port scanner to make sure the server port is close.  
If it detects open ports, an error will be output and the exit status will end with 1 or higher.  
Optional: slack notification if it get errors.

## Usage
### go run command.
serverlist.txtを標準入力から渡して複数サーバをスキャンします。  
serverlist.txtは1行1サーバの改行区切りで、コメントは`#`を用います。  
DNSで名前解決できないサーバの場合はスキャンをスキップします（エラーにはなりません）。
```
go run aite9.go -tcp 22,25,3306  < serverlist.txt
```

## options
### silent mode option
`-mode silent` no message to stdout. only output error messages to stderr.
```
go run aite9.go -tcp 22,25,3306 -mode silent < serverlist.txt
```

## Results
Return status code 0 if there is no problem.  
Return status code 1 or higher with error message if there there are problems.

## Slack Notification
Set slack webhook settings on OS env, 
NSchecker sends error message to the slack channel when detecting errors or DNS record changing.

```cassandraql
export SLACK_WEBHOOK_URL="webhook url"
export SLACK_FREE_TEXT="<!channel> from AWS lambda" #optional
export SLACK_USERNAME="your user" #optional
export SLACK_CHANNEL="your channel" #optional
export SLACK_ICON_EMOJI=":smile:" #optional
export SLACK_ICON_URL="icon url" #optional
``` 
