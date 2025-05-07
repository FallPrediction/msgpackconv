# msgpackconv
這是一個小實作，將 message pack 轉為 JSON，或將 JSON 轉為 message

## Quick start
```go
msg := msgpack.FromJSON([]byte(`{"str": "a"}`)) // 讀取 JSON 字串
fmt.Printf("% x\n", msg) // 轉為 message pack bytes 81 a3 73 74 72 a1 61

json := msgpack.ToJSON(msg) // 讀取 message pack bytes
fmt.Println(string(json)) // 轉為 JSON {"str":"a"}

// FromJSON 和 ToJSON 遇到 invalid input 皆輸出 empty byte slice
```

## JSON 轉 message pack
解析一個結構未知的 JSON 為一個 empty interface 變數，然後因為 interface value 保存它底層的具體類型和值，所以可以利用 type switch 存取它的底層資料類型和值，轉換成message pack 相應的資料類型、長度和資料本身

舉例：

解析一個 JSON 為變數 obj（empty interface），利用 type switch 檢查出底層 type 為 string

然後計算 value 長度，>=32 且 < 2 的 8 次方，對應的 message pack type 為 str 8
- 第一個 byte 為 `0xd9`
- 第二個 byte 為資料長度
- 其餘的 byte 為資料本身

然後將這些 `[]byte` 放入最終答案裡

如果解析出的 type 是 slice 或 map，則利用遞迴解析出每個元素或 key-value pair

## message pack 轉 JSON
依序讀取輸入 message pack 的 bytes，找到特定資料類型的 first byte 時，解析出資料長度和資料本身，再利用反射設定 empty interface 的底層 value。全部讀取後，將該 interface 轉為 JSON

舉例：

設定一個 empty interface 變數，開始讀取 message pack
- 第 0 個 byte，讀取到 `0xd9`，表示資料為 str8
- 讀取下一個 byte 取得資料長度
- 依照長度讀取後面的 byte 取得資料本身

然後將這個 interface 的 value 設為該 string

如果讀取的資料類型為 array 或 map，則迴圈讀取每個元素或 key-value pair

## 參考
- [MessagePack 規範](https://github.com/msgpack/msgpack/blob/master/spec.md)
- [JSON and Go - The Go Programming Language](https://go.dev/blog/json)
