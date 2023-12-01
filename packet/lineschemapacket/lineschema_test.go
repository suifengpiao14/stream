package lineschemapacket_test

import (
	"fmt"
	"testing"

	"github.com/tidwall/gjson"
)

/*
*
{pageInfo:{pageIndex:input.pageIndex,pageSize:input.pageSize,total:PaginateTotalOut},items:{key:PaginateOut.#.key,label:PaginateOut.#.label,title:PaginateOut.#.title,deletedAt:PaginateOut.#.deleted_at,description:PaginateOut.#.description,id:PaginateOut.#.id,thumb:PaginateOut.#.thumb,updatedAt:PaginateOut.#.updated_at,content:PaginateOut.#.content,createdAt:PaginateOut.#.created_at,icon:PaginateOut.#.icon}|@group}
*
*/
func TestPath(t *testing.T) {
	input := `{"code":200,"message":"ok","items":[{"id":1,"title":"test1"},{"id":2,"title":"test2"}],"pagination":{"index":0,"size":10,"total":100}}`
	pathMap := "{code:code.@tostring,message:message.@tostring,items.#.id:items.#.id.@tostring,items.#.title:items.#.title.@tostring,}"
	//pathMap = `{code:code.@tostring,message:message.@tostring,pagination:{index:pagination.index.@tostring,size:pagination.size.@tostring,total:pagination.total.@tostring},items:{id:items.#.id.@tostring,title:items.#.title.@tostring}|@group}`
	outputStr := gjson.GetBytes([]byte(input), pathMap).String()
	fmt.Println(outputStr)
}
