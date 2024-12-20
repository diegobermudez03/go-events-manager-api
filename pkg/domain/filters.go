package domain

import "github.com/google/uuid"

// FILTER OBJECTS
type ParticipationFilters struct{
	UserId 		*uuid.UUID 
	RoleName 	*string 
	EventId 	*uuid.UUID
	Offset 		*int 
	Limit 		*int 
}

type ParticipationFilter func(filter *ParticipationFilters)

func ParticipationUserIdFilter(userId *uuid.UUID) ParticipationFilter{
	return func(filter *ParticipationFilters) {
		filter.UserId = userId
	}
}

func ParticipationRoleFilter(roleName *string) ParticipationFilter{
	return func(filter *ParticipationFilters){
		filter.RoleName = roleName
	}
}

func ParticipationEventIdFilter(eventId *uuid.UUID) ParticipationFilter{
	return func(filter *ParticipationFilters) {
		filter.EventId = eventId
	}
}

func ParticipationOffsetFilter(offset *int) ParticipationFilter{
	return func(filter *ParticipationFilters) {
		filter.Offset = offset
	}
}

func ParticipationLimitFilter(limit *int) ParticipationFilter{
	return func(filter *ParticipationFilters) {
		filter.Limit = limit
	}
}