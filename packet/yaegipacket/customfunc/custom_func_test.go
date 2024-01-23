package customfunc_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/httpraw"
	"github.com/suifengpiao14/stream/packet/yaegipacket/customfunc"
)

func TestGetSetJson(t *testing.T) {
	body := `{"_head":{"_version":"0.01","_msgType":"request","_timestamps":"1705903583","_interface":"dispatch_get_order_list","_remark":""},"_param":{"ouid":"","ouname":"","operaterId":"","operaterName":"","creatorId":"","creator":"","dispatchStatus":"","orderChannel":"","orderId":"","orderNumber":"2401221118100134","orderTag":"","orderType":"","userPhone":"","startTime":"","endTime":"","cityId":"","pageSize":"","pageNum":"3","SFFlag":"","desc":"导出操作"}}`
	rDTO := httpraw.RequestDTO{
		Body: body,
	}
	b, err := json.Marshal(rDTO)
	require.NoError(t, err)
	input := string(b)
	gjsonpath := `body.@fromstr._param.pageNum`

	params, err := customfunc.GetsetJson(input, gjsonpath, func(oldValue string) (newValue string, err error) {
		index := cast.ToInt(oldValue)
		// if index == 0 {
		// 	return cast.ToString(index), nil
		// }
		index++
		return cast.ToString(index), nil
	})
	require.NoError(t, err)
	fmt.Println(params)
}
