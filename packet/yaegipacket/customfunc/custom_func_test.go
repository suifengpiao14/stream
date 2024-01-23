package customfunc_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/httpraw"
	"github.com/suifengpiao14/stream/packet/yaegipacket/customfunc"
	"github.com/tidwall/gjson"
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

func TestWalkArrayMap(t *testing.T) {
	input := `{"_data":{"_data":{"list":[{"addAmount":"0","amount":"240000","city":"重庆市","cityId":"5001","clerkId":"1272410","clerkName":"谢泽云","clerkRole":"S1","closeReason":"","creditFlag":"0","dispatchId":"1722883","dispatchStatus":"7","dispatchTime":"2024-01-22 13:57:49","orderChannel":"17","orderId":"20762164","orderNum":"2401221118100134","orderTime":"2024-01-22 13:47:39","phoneModel":"魅族 20 INFINITY 无界版","refundRemark":"","storeId":"1287109","storeName":"","updateTime":"2024-01-22 14:01:34","userAddr":"重庆重庆市九龙坡区谢家湾街道黄杨路24号大鼎世纪滨江3栋1504","userExpectTime":"2024-01-27 09:00:00","userName":"黄献","userPhone":"19122085667"}],"listLen":"1","number":"1"},"_errCode":"0","_errStr":"SUCCEED","_ret":"0"},"_head":{"_interface":"dispatch_get_order_list","_msgType":"response","_remark":"","_timestamps":"1705909960","_version":"0.01"}}`

	statusMap := `{"1":"待派单","2":"已派单","3":"已上门","4":"已到店","5":"已回收","6":"上门后关闭","7":"上门前关闭","8":"丰修待派单","9":"关闭中","10":"待提交新机单","11":"待发货"}`
	statusMapping := gjson.Parse(statusMap)
	storeNameMapping := gjson.Parse(`{"":"待指派门店"}`)
	output, err := customfunc.WalkArrayMap(input, "_data._data.list", func(row map[string]any) (newRow map[string]any) {
		val := cast.ToString(row["storeName"])
		row["storeName"] = customfunc.FiledMapping(val, storeNameMapping, val)
		row["dispatchStatus"] = statusMapping.Get(cast.ToString(row["dispatchStatus"])).String()
		row["amount"] = customfunc.Fen2yuan(row["amount"])
		row["addAmount"] = customfunc.Fen2yuan(row["addAmount"])
		row["userPhone"] = customfunc.MaskFunc(cast.ToString(row["userPhone"]), "****", 3, 4)
		return row
	})
	require.NoError(t, err)
	fmt.Println(output)
}
