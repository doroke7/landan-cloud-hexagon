package pkgGin

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type Request struct {
	Context *gin.Context
}

func (oSelf *Request) Bind(sString string, oValue any) error {

	if sString == "filters" {
		sFilters := oSelf.Context.Query("filters")
		return json.Unmarshal([]byte(sFilters), oValue)
	}

	if sString == "sorters" {
		sSorters := oSelf.Context.Query("sorters")
		return json.Unmarshal([]byte(sSorters), oValue)
	}

	if sString == "pagination" {
		sPagination := oSelf.Context.Query("pagination")
		return json.Unmarshal([]byte(sPagination), oValue)
	}

	if sString == "value" {
		sValue := oSelf.Context.PostForm("value")
		return json.Unmarshal([]byte(sValue), oValue)
	}

	return nil
}
