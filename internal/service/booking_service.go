package service

import (
	"errors"
	"fmt"
	log "log/slog"
	"test-backend-1-d1ma11/internal/entity"
	"test-backend-1-d1ma11/internal/repository"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BookingServiceImpl struct {
	bookingRepo       repository.BookingRepository
	slotRepo          repository.SlotRepository
	conferenceService ConferenceService
}

func NewBookingService(bookingRepo repository.BookingRepository, slotRepo repository.SlotRepository, conferenceService ConferenceService) *BookingServiceImpl {
	return &BookingServiceImpl{bookingRepo: bookingRepo, slotRepo: slotRepo, conferenceService: conferenceService}
}

func (s *BookingServiceImpl) Create(booking *entity.Booking, createConferenceLink bool) error {
	log.Info("Starting creation booking for slot")
	err := s.bookingRepo.(*repository.BookingRepositoryImpl).DB.Transaction(func(tx *gorm.DB) error {
		var slot entity.Slot
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&slot, "id = ?", booking.SlotID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Info(fmt.Sprintf("There is no slot with id=%v for booking", booking.SlotID))
				return NewError(ErrorType.SlotNotFound, "slot not found")
			}
			log.Info(fmt.Sprintf("Failed to find slot with id=%v. Error: %v", booking.SlotID, err))
			return err
		}

		if slot.StartAt.Before(time.Now().UTC()) {
			log.Info("Slot start is in the past, but must be in the future")
			return NewError(ErrorType.InvalidRequest, "cannot book a slot in the past")
		}

		var existing entity.Booking
		if err := tx.Where("slot_id = ? AND status = ?", booking.SlotID, "active").First(&existing).Error; err == nil {
			log.Info(fmt.Sprintf("Slot with id=%v is already booked", booking.SlotID))
			return NewError(ErrorType.SlotAlreadyBooked, "slot is already booked")
		}

		booking.Status = "active"
		if err := tx.Create(booking).Error; err != nil {
			return err
		}

		var conferenceLink *string
		if createConferenceLink {
			link, err := s.conferenceService.CreateConferenceLink(booking.ID)
			if err != nil {
				log.Warn(fmt.Sprintf("Failed to create conference link: %v", err))
				return nil
			}
			conferenceLink = &link
			booking.ConferenceLink = conferenceLink

			if err = tx.Save(booking).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		var appErr *ErrorResponse
		if errors.As(err, &appErr) {
			return err
		}
		log.Info(fmt.Sprintf("Failed to create booking for slot with id=%v. Error: %v", booking.SlotID, err))
		return NewInternalError("failed to create booking")
	}

	return nil
}

func (s *BookingServiceImpl) ListByUserId(bookings *[]entity.Booking, userId string) error {
	if err := s.bookingRepo.GetByUserID(bookings, userId); err != nil {
		return NewInternalError("failed to get user's bookings")
	}
	return nil
}

func (s *BookingServiceImpl) ListByAdminId(bookings *[]entity.Booking, page, pageSize int, total *int64) error {
	if err := s.bookingRepo.GetAllWithPagination(bookings, page, pageSize, total); err != nil {
		return err
	}
	return nil
}

func (s *BookingServiceImpl) CancelBooking(booking *entity.Booking, bookingId, userId string) error {
	log.Info(fmt.Sprintf("%v", booking))
	if err := s.bookingRepo.GetById(booking, bookingId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NewError(ErrorType.BookingNotFound, "booking not found")
		}
		return err
	}
	log.Info(fmt.Sprintf("%v", booking))

	if booking.UserID != userId {
		return NewError(ErrorType.Forbidden, "cannot cancel another user's booking")
	}

	if booking.Status == "cancelled" {
		return nil
	}

	booking.Status = "cancelled"
	if err := s.bookingRepo.Save(booking); err != nil {
		return NewInternalError("failed to cancel booking")
	}

	return nil
}
