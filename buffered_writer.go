package cache

import (
	"bytes"
	"fmt"
	"io"

	gds "github.com/zr-hebo/gdstructure"
)

const (
	bufferNum = 2
)

// BufferedWriter BufferedWriter
type BufferedWriter struct {
	writer       io.Writer
	blockPool    *gds.Queue
	hotBuff      *bytes.Buffer
	bufferQueue  chan *bytes.Buffer
	finishWaiter chan struct{}
	errWaiter    chan error
	bufferSize   int
}

// NewBufferedWriter NewBufferedWriter
func NewBufferedWriter(
	writer io.Writer, bufferSize int) (bfw *BufferedWriter, err error) {
	if writer == nil {
		err = fmt.Errorf("传入的Writer不能为Nil，请检查确认！")
		return
	}

	bfw = new(BufferedWriter)
	bfw.blockPool = gds.NewQueue()
	bfw.bufferSize = bufferSize
	bfw.bufferQueue = make(chan *bytes.Buffer, bufferNum)
	bfw.finishWaiter = make(chan struct{}, 1)
	bfw.errWaiter = make(chan error, 1)
	bfw.writer = writer

	go bfw.writeFile()

	return
}

func (p *BufferedWriter) writeFile() {
	defer func() {
		if panicRecover := recover(); panicRecover != nil {
			p.errWaiter <- fmt.Errorf(
				"写入文件的时候panic <-- %v", panicRecover)
		}
	}()

	var bufferData *bytes.Buffer
	for {
		// 从写缓存池提取可写内容
		select {
		case err := <-p.errWaiter:
			p.errWaiter <- err
			return
		default:
			bufferData = <-p.bufferQueue
		}

		// 设置写完成
		if bufferData == nil {
			p.finishWaiter <- struct{}{}
			return
		}

		// 写文件
		if _, err := p.writer.Write(bufferData.Bytes()); err != nil {
			var writeBuffErr error
			err = fmt.Errorf("写入文件失败 <-- %s", err.Error())
			select {
			case p.errWaiter <- err:
			case writeBuffErr = <-p.errWaiter:
				err = fmt.Errorf("%s; %s", err.Error(), writeBuffErr.Error())
				p.errWaiter <- err
			}
			return
		}

		// 将使用过的内存块缓存起来供以后使用
		bufferData.Reset()
		p.blockPool.Enqueue(bufferData)
	}
}

func (p *BufferedWriter) claimSpace() {
	// 查看是否由空闲可用的内存块可以使用
	// 没有的话，申请新的内存空间
	block := p.blockPool.Dequeue()
	if block != nil {
		p.hotBuff = block.(*bytes.Buffer)

	} else {
		bufferBytes := make([]byte, 0, p.bufferSize)
		p.hotBuff = bytes.NewBuffer(bufferBytes)
	}
}

// Write 将buffer写入缓存
func (p *BufferedWriter) Write(data []byte) (err error) {
	defer func() {
		var writeFileErr error
		if err != nil {
			err = fmt.Errorf("写入缓存失败 <-- %s", err.Error())
		}

		select {
		case writeFileErr = <-p.errWaiter:
			if err != nil {
				err = fmt.Errorf("%s; %s", writeFileErr.Error(), err.Error())
			}
			p.errWaiter <- err

		default:
			if err != nil {
				p.errWaiter <- err
			}
		}
	}()

	// 判断当前写入缓存是否为空
	if p.hotBuff == nil {
		p.claimSpace()
	}

	// 判断当前缓存是否可以继续写入
	if p.hotBuff.Len()+len(data) > p.hotBuff.Cap() {
		// 将原来缓存的数据提交写入文件
		p.bufferQueue <- p.hotBuff
		p.claimSpace()

		// 判断新申请写入数据本身大过缓存大小
		// 如果写入数据过大，直接将数据提交写入文件
		if len(data) >= p.hotBuff.Cap() {
			p.bufferQueue <- bytes.NewBuffer(data)

		} else {
			p.hotBuff.Write(data)
		}

	} else {
		_, err = p.hotBuff.Write(data)
	}

	return
}

// commitAllWrite 将当前缓存剩余的数据提交写入文件
func (p *BufferedWriter) commitAllWrite() {
	p.bufferQueue <- p.hotBuff
	close(p.bufferQueue)
	return
}

// WaitClean 等待所有缓存写入完成或者失败
func (p *BufferedWriter) WaitClean() (err error) {
	p.commitAllWrite()

	select {
	case <-p.finishWaiter:
	case err = <-p.errWaiter:
	}
	return
}
