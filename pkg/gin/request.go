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
		oErr := json.Unmarshal([]byte(sFilters), oValue)

		return oErr
	}

	if sString == "sorters" {
		sSorters := oSelf.Context.Query("sorters")
		oErr := json.Unmarshal([]byte(sSorters), oValue)

		return oErr
	}

	if sString == "pagination" {
		sPagination := oSelf.Context.Query("pagination")
		oErr := json.Unmarshal([]byte(sPagination), oValue)

		return oErr
	}

	if sString == "variable" {
		sVariable := oSelf.Context.PostForm("variable")
		oErr := json.Unmarshal([]byte(sVariable), oValue)

		return oErr
	}

	return nil
}
