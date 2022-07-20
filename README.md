# go_cobra
Go言語とCobraを使ったコマンドラインツールを自作するリポジトリ

## get started

- cobraを入手
`go get github.com/spf13/cobra/cobra`  
もしくは
`go get github.com/spf13/cobra@latest`

- cobraコマンドが使えない時
`$HOME/go/bin`配下にcobraが設置されている時、エイリアスを設定  
zshrcを編集  
`alias cobra="$HOME/go/bin/cobra-cli"`

- コマンド作成
`cobra add [command_name]`

- コマンド実行
`go run main.go [command_name]`

## Example
hello.goを実行  
`go main.go hello -l ja Taro Jiro`  
実行結果  
`Taroさん, Jiroさん, こんにちは！`  

グローバルフラッグについてはroot.goのinitで定義

