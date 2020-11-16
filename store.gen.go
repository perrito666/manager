// Code generated by oto; DO NOT EDIT.

package main

import ()

type Store interface {
	CreateEventSlot(EventSlot) error
	GetEventSlot(EventSlot) (EventSlot, error)
	GetEventSlotByID(uint32) (EventSlot, error)

	UpdateEventSlot(EventSlot) error
	DeleteEventSlot(EventSlot) error
	DeleteEventSlotByID(uint32) (EventSlot, error)

	CreateEvent(Event) error
	GetEvent(Event) (Event, error)
	GetEventByID(uint32) (Event, error)
	GetEventSlots(uint32) ([]Slots, error)
	UpdateEvent(Event) error
	DeleteEvent(Event) error
	DeleteEventByID(uint32) (Event, error)

	CreateConference(Conference) error
	GetConference(Conference) (Conference, error)
	GetConferenceByID(uint32) (Conference, error)
	GetConferenceEvents(uint32) ([]Events, error)
	UpdateConference(Conference) error
	DeleteConference(Conference) error
	DeleteConferenceByID(uint32) (Conference, error)
}
