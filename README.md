一个golang实现的处理缓存的包  

### bufferedwriter
 写文件时，为了提高性能，写的内容先放在缓存中，等缓存满了的时候再写入到文件。
#### Install

Use go get to install this package:<br>
$ go get github.com/zr-hebo/util-cache

#### Usage

```golang
import (
	cu "github.com/zr-hebo/util-cache"
)

// 创建带缓存的Writer，file是打开的文件，bufferSize是缓存的大小
writer, err := cu.NewBufferedWriter(file, bufferSize)
if err != nil {
	return
}

// 将byte数组写入writer	
if err = writer.Write(rowData); err != nil {
	return
}

// 等待缓存中的数据全部写入文件
if err = writer.WaitClean(); err != nil {
	return
}
```
