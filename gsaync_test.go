package gasync

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go/types"
	"testing"
	"time"
)

func Test_P0(t *testing.T) {
	msg := "hello ~ I am goAsync ~"
	f1 := GoAsync[string](func() (string, error) {
		time.Sleep(time.Second * 5)
		return msg, nil
	})
	GoAsync[types.Nil](func() (types.Nil, error) {
		i := 0
		for true {
			fmt.Println(i)
			i++
			time.Sleep(time.Second)
		}
		return types.Nil{}, nil
	})
	_, done, _ := f1.GetNow()
	assert.False(t, done)
	get, err := f1.Get()
	fmt.Println(get)
	assert.Nil(t, err)
	assert.Equal(t, msg, get)

}
