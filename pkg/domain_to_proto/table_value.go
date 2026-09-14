package domainToProto

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	domain "example/internal/domain"

	pbResource "example/pb/resource"
)

func TableValue(oTable *domain.TableVariable) *pbResource.TableValue {
	oValue := &pbResource.TableValue{
		No:          oTable.No,
		Key:         oTable.Key,
		Description: oTable.Description,
	}

	if oTable.GameId != nil {
		iGameId := uint64(*oTable.GameId)
		oValue.GameId = &iGameId
	}

	if oTable.State != nil {
		iState := uint64(*oTable.State)
		oValue.State = &iState
	}

	if oTable.Result != nil {
		sResult := *oTable.Result
		oValue.Result = &sResult
	}

	if oTable.StartedAt != nil {
		oValue.StartedAt = timestamppb.New(*oTable.StartedAt)
	}

	if oTable.EndedAt != nil {
		oValue.EndedAt = timestamppb.New(*oTable.EndedAt)
	}

	return oValue
}
