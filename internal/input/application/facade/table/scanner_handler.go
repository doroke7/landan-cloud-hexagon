package applicationFacadeTable

import (
	inputApplicationFacade "example/internal/input/application/facade"
	pbFacadeTable "example/pb/facade/table"
)

type ScannerHandler struct {
	pbFacadeTable.UnimplementedScannerControllerServer
	*inputApplicationFacade.AbstractHandler
}

func NewScannerHandler(oAbstractHandler *inputApplicationFacade.AbstractHandler) *ScannerHandler {
	return &ScannerHandler{
		AbstractHandler: oAbstractHandler,
	}
}
