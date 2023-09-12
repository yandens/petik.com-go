package helpers

import "gorm.io/gorm"

func IsUniqueSeatNumber(seatNumbers []string) bool {
  // make a map of string to bool
  unique := make(map[string]bool)

  // loop through the seatNumbers
  for _, seatNumber := range seatNumbers {
    if unique[seatNumber] {
      return false
    }
    unique[seatNumber] = true
  }

  return true
}

func IsSeatNumberValid(seatNumbers []string) bool {
  // make a map of string to bool
  valid := make(map[string]bool)

  // loop through the seatNumbers
  for _, seatNumber := range seatNumbers {
    if seatNumber == "" {
      return false
    }
    valid[seatNumber] = true
  }

  return true
}

func IsSeatNumberAvailable(seatNumbers []string, flightID uint, db *gorm.DB) bool {
  // find the reserved seat numbers
  var reservedSeatNumbers []string
  db.Table("booking_details").Select("seat_number").Joins("JOIN bookings ON bookings.id = booking_details.booking_id").Where("bookings.flight_id = ?", flightID).Scan(&reservedSeatNumbers)

  // loop through the seatNumbers and reservedSeatNumbers
  for _, seat := range seatNumbers {
    for _, reservedSeat := range reservedSeatNumbers {
      if seat == reservedSeat {
        return false
      }
    }
  }

  return true
}
