package pkg

func PaginationToLimit(oPagination *Pagination) *Limit {
	iSize := uint(10)
	iPage := uint(1)

	if oPagination != nil && oPagination.Size != nil && *oPagination.Size != 0 {
		iSize = *oPagination.Size
	}

	if oPagination != nil && oPagination.Page != nil && *oPagination.Page != 0 {
		iPage = *oPagination.Page
	}

	iOffset := (iPage - 1) * iSize

	oLimit := &Limit{
		Offset: &iOffset,
		Count:  &iSize,
	}

	return oLimit
}
