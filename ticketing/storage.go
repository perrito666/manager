package ticketing

import (
	"fmt"

	"github.com/ShiftLeftSecurity/gaum/db/chain"
	"github.com/ShiftLeftSecurity/gaum/db/connection"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	uuid "github.com/satori/go.uuid"
)

type SQLStorage struct {
	conn connection.DB
}

// AtomicOperation begins a transaction and returns commit and rollback functions along with a new
// SQLStorage wrapping the tx
func (s *SQLStorage) AtomicOperation() (func() error, func() error, PurchaseStore, error) {
	tx, err := s.conn.BeginTransaction()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("beginning transation: %w", err)
	}
	return tx.CommitTransaction, tx.RollbackTransaction, &SQLStorage{conn: tx}, nil
}

const (
	ticketIDUniqueConstraint = "ticket_id_is_unique"
	tableSlotClaims          = "slot_claim"
)

// CreateSlotClaim saves a slot claim and returns it with the populated ID
func (s *SQLStorage) CreateSlotClaim(slotClaim *SlotClaim) (*SlotClaim, error) {
	var err error
	// FIXME: Add a check for capacity not exceeded on event.
	for i := 0; i < 3; i++ {
		q := chain.New(s.conn)
		results := []SlotClaim{}
		err = q.Insert(map[string]interface{}{
			"ticket_id":     slotClaim.TicketID,
			"redeemed":      slotClaim.Redeemed,
			"event_slot_id": slotClaim.EventSlot.ID,
		}).
			Table(tableSlotClaims).
			Returning("id, ticket_id, redeemed").Fetch(&results)

		if err != nil {
			// there is no SQL for "on error change the inserting statement", only to change
			// the existing one.
			if err, ok := err.(pgx.PgError); ok {
				if pgerrcode.IsIntegrityConstraintViolation(err.Code) {
					// if this clashes again entropy might be broken, check if not 2020
					slotClaim.TicketID = uuid.NewV4().String()
					continue
				}
			}
			return nil, fmt.Errorf("saving slot claim: %w", err)
		}
		if len(results) == 0 {
			return nil, nil
		}
		results[0].EventSlot = slotClaim.EventSlot
		return &results[0], nil
	}
	return nil, fmt.Errorf("failed to insert: %w", err)
}

const (
	tableAttendee               = "attendee"
	tableAttendeeSlotClaims     = "attendee_to_slot_claims"
	slotClaimIDUniqueConstraint = "slot_claim_id_is_unique"
)

// UpdateAttendee saves the passed attendee attributes on top of the existing one.
func (s *SQLStorage) UpdateAttendee(attendee *Attendee) (*Attendee, error) {
	q := chain.New(s.conn)
	rows, err := q.UpdateMap(map[string]interface{}{
		"email":        attendee.Email,
		"coc_accepted": attendee.CoCAccepted,
	}).
		AndWhere("id = ?", attendee.ID).ExecResult()
	if err != nil {
		return nil, fmt.Errorf("updating attendee: %w", err)
	}
	if rows == 0 { // attendee does not exist
		return nil, nil
	}

	for i := range attendee.Claims {
		c := attendee.Claims[i]
		err := q.Insert(map[string]interface{}{
			"attendee_id":   attendee.ID,
			"slot_claim_id": c.Redeemed,
		}).Table(tableAttendeeSlotClaims).
			OnConflict(func(c *chain.OnConflict) {
				// This claim was someone else's, this might be the result of transfering.
				c.OnConstraint(slotClaimIDUniqueConstraint).
					DoUpdate().
					Set("attendee_id", attendee.ID)
			}).Exec()
		if err != nil {
			return nil, fmt.Errorf("updating attendee claims: %w", err)
		}

	}
	return attendee, nil
}

func (s *SQLStorage) CreateClaimPayment(_ *ClaimPayment) (*ClaimPayment, error) {
	panic("not implemented") // TODO: Implement
}

func (s *SQLStorage) UpdateClaimPayment(_ *ClaimPayment) (*ClaimPayment, error) {
	panic("not implemented") // TODO: Implement
}

func (s *SQLStorage) ChangeSlotClaimOwner(_ []SlotClaim, _ *Attendee, _ *Attendee) (*Attendee, *Attendee, error) {
	panic("not implemented") // TODO: Implement
}
