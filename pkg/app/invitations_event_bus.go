package app

import (
	"sync"

	"github.com/diegobermudez03/go-events-manager-api/pkg/domain"
	"github.com/google/uuid"
)

type InvitationsEventBus struct {
	Suscribers map[uuid.UUID] []chan domain.InvitationEvent
	Mutex sync.RWMutex
}


func NewInvitationsEventBus() domain.InvitationsEventBus{
	return &InvitationsEventBus{
		Suscribers: map[uuid.UUID][] chan domain.InvitationEvent{},
		Mutex: sync.RWMutex{},
	}
}

func (b *InvitationsEventBus) Publish(eventId uuid.UUID, event domain.InvitationEvent){
	go func(){
		b.Mutex.RLock()
		defer b.Mutex.RUnlock()
		suscribers, ok := b.Suscribers[eventId]
		if !ok{
			return
		}
		for _, sus := range suscribers{
			select{
			case sus <- event:
			default:
			}
		}
	}()
}

func (b *InvitationsEventBus) Suscribe(eventId uuid.UUID) <-chan domain.InvitationEvent{
	channel := make(chan domain.InvitationEvent)
	b.Mutex.Lock()
	defer b.Mutex.Unlock()
	if suscribers, ok := b.Suscribers[eventId]; ok{
		b.Suscribers[eventId] = append(suscribers, channel)
	}else{
		b.Suscribers[eventId] = []chan domain.InvitationEvent{channel}
	}
	return channel
}


func (b *InvitationsEventBus) Unsuscribe(eventId uuid.UUID, channel <-chan domain.InvitationEvent){
	b.Mutex.Lock()
	suscribers := b.Suscribers[eventId]
	for i, sus := range suscribers{
		if sus == channel{
			b.Suscribers[eventId] = append(b.Suscribers[eventId][:i], b.Suscribers[eventId][i+1:]...)
			break
		}
	}
	b.Mutex.Unlock()
}