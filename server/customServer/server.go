package customserver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	pb "github.com/pkkp86nitd/train_ticket_app/proto"
)

const (
	numSections     = 2
	seatsPerSection = 10
)

type TrainTicketServer struct {
	pb.UnimplementedTrainTicketServiceServer
	tickets      map[string]*pb.Receipt
	mu           sync.Mutex
	seatOccupied map[string]int
}

var (
	instance *TrainTicketServer
	once     sync.Once
)

// GetCustomServerInstance returns a singleton instance of TrainTicketServer
func GetCustomServerInstance() *TrainTicketServer {
	once.Do(func() {
		instance = &TrainTicketServer{
			tickets:      make(map[string]*pb.Receipt),
			seatOccupied: make(map[string]int),
		}
	})
	return instance
}

// NewTrainTicketServer creates a new instance of TrainTicketServer
func NewTrainTicketServer() *TrainTicketServer {
	return &TrainTicketServer{
		tickets:      make(map[string]*pb.Receipt),
		seatOccupied: make(map[string]int),
	}
}

func (s *TrainTicketServer) SubmitPurchase(ctx context.Context, req *pb.PurchaseRequest) (*pb.Receipt, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	log.Printf("Received SubmitPurchase request for user: %s", req.UserEmail)

	for _, section := range []string{"A", "B"} {
		for i := 0; i < seatsPerSection; i++ {
			seatSection := fmt.Sprintf("%s%d", section, i+1)
			if !s.IsSeatOccupied(seatSection) {
				price := 20.0 // Fixed price

				receipt := &pb.Receipt{
					From:          req.From,
					To:            req.To,
					UserFirstName: req.UserFirstName,
					UserLastName:  req.UserLastName,
					UserEmail:     req.UserEmail,
					SeatSection:   section,
					SeatNumber:    int32(i + 1),
					PricePaid:     price,
				}


				time.Sleep(1 * time.Second)

				// Update the seat occupancy
				s.seatOccupied[seatSection]++
				s.tickets[req.UserEmail] = receipt

				log.Printf("Purchase successful for user: %s, Seat Section: %s, Seat Number: %d", req.UserEmail, section, i+1)
				return receipt, nil
			}
		}
	}

	log.Printf("Purchase failed for user: %s, all sections are full", req.UserEmail)
	return nil, errors.New("all sections are full")
}

func (s *TrainTicketServer) ViewReceipt(ctx context.Context, req *pb.ViewReceiptRequest) (*pb.Receipt, error) {
	log.Printf("Received ViewReceipt request for user: %s", req.UserEmail)

	if receipt, ok := s.tickets[req.UserEmail]; ok {
	
		time.Sleep(1 * time.Second)

		log.Printf("Returning receipt for user: %s", req.UserEmail)
		return receipt, nil
	}

	log.Printf("Receipt not found for user: %s", req.UserEmail)
	return nil, fmt.Errorf("receipt not found for user email: %s", req.UserEmail)
}

func (s *TrainTicketServer) ViewUsersBySection(ctx context.Context, req *pb.SectionRequest) (*pb.UsersBySectionResponse, error) {
	log.Printf("Received ViewUsersBySection request for section: %s", req.Section)

	usersBySection := make([]*pb.UserSeat, 0)
	for userEmail, receipt := range s.tickets {
		//fmt.Printf("Checking user: UserEmail=%s, ReceiptSection=%s, RequestedSection=%s, comparison=%t\n", userEmail, receipt.SeatSection, req.Section, strings.TrimSpace(receipt.SeatSection) == strings.TrimSpace(req.Section))
		if receipt.SeatSection == req.Section {
			fmt.Printf("Adding user to result: UserEmail=%s, SeatSection=%s, SeatNumber=%d, TicketSection=%s\n",
				userEmail, receipt.SeatSection, receipt.SeatNumber, receipt.SeatSection)
			usersBySection = append(usersBySection, &pb.UserSeat{
				UserEmail:   userEmail,
				SeatSection: receipt.SeatSection,
				SeatNumber:  receipt.SeatNumber,
			})
		}
	}

	log.Printf("Returning %d users for section: %s", len(usersBySection), req.Section)

	time.Sleep(1 * time.Second)

	return &pb.UsersBySectionResponse{UserSeats: usersBySection}, nil
}

func (s *TrainTicketServer) RemoveUser(ctx context.Context, req *pb.UserRequest) (*pb.RemoveUserResponse, error) {
	
	log.Printf("Received RemoveUser request for user: %s", req.UserEmail)

	// Implement logic to remove a user
	if receipt, ok := s.tickets[req.UserEmail]; ok {
		// Update the seat occupancy
		log.Printf("Receipt details: %s", receipt)
		s.seatOccupied[fmt.Sprintf("%s%d", receipt.SeatSection, receipt.SeatNumber)]--
		delete(s.tickets, req.UserEmail)

		// Introduce a 1-second delay
		time.Sleep(1 * time.Second)

		log.Printf("User removed successfully: %s", req.UserEmail)
		return &pb.RemoveUserResponse{Success: true}, nil
	}

	log.Printf("User not found for removal: %s", req.UserEmail)
	return &pb.RemoveUserResponse{Success: false}, fmt.Errorf("user not found for email: %s", req.UserEmail)
}

func (s *TrainTicketServer) ModifySeat(ctx context.Context, req *pb.ModifySeatRequest) (*pb.ModifySeatResponse, error) {
	log.Printf("Received ModifySeat request for user: %s", req.UserEmail)

	s.mu.Lock()
	defer s.mu.Unlock()

	if receipt, ok := s.tickets[req.UserEmail]; ok {
		// Attempt to update the seat to the specified new section
		newSeatSection := fmt.Sprintf("%s%d", req.NewSection, receipt.SeatNumber)
		if !s.IsSeatOccupied(newSeatSection) {
			// Seat in the new section is available, update the seat
			s.seatOccupied[fmt.Sprintf("%s%d", receipt.SeatSection, receipt.SeatNumber)]--
			s.seatOccupied[newSeatSection]++
			receipt.SeatSection = req.NewSection

			time.Sleep(1 * time.Second)

			log.Printf("Seat modified successfully for user: %s, New Section: %s", req.UserEmail, req.NewSection)
			return &pb.ModifySeatResponse{Success: true}, nil
		}

		// If the seat in the new section is not available, find an available seat in another section
		for i := 0; i < numSections; i++ {
			alternativeSeatSection := fmt.Sprintf("%s%d", req.NewSection, i+1)
			if !s.IsSeatOccupied(alternativeSeatSection) {
				// Found an available seat in another section, update the seat
				s.seatOccupied[fmt.Sprintf("%s%d", receipt.SeatSection, receipt.SeatNumber)]--
				s.seatOccupied[alternativeSeatSection]++
				receipt.SeatSection = alternativeSeatSection

				time.Sleep(1 * time.Second)

				log.Printf("Seat modified successfully for user: %s, New Section: %s", req.UserEmail, alternativeSeatSection)
				return &pb.ModifySeatResponse{Success: true}, nil
			}
		}

		log.Printf("No available seats found for user: %s in the new section or alternative section", req.UserEmail)
		return &pb.ModifySeatResponse{Success: false}, errors.New("no available seats found")
	}

	log.Printf("User not found for seat modification: %s", req.UserEmail)
	return &pb.ModifySeatResponse{Success: false}, fmt.Errorf("user not found for email: %s", req.UserEmail)
}

func (s *TrainTicketServer) IsSeatOccupied(seatSection string) bool {

	return s.seatOccupied[seatSection] > 0
}

func (s *TrainTicketServer) GetTickets() map[string]*pb.Receipt {
	return s.tickets
}


func (s *TrainTicketServer) GetSeatOccupied() map[string]int {
	return s.seatOccupied
}
