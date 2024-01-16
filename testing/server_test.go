package testing

import (
    "context"
    "fmt"
    "testing"

    customserver "github.com/pkkp86nitd/train_ticket_app/server/customServer"
    "github.com/stretchr/testify/assert"
    pb "github.com/pkkp86nitd/train_ticket_app/proto"
)

func TestSubmitPurchase(t *testing.T) {

	fmt.Printf("------ TestSubmitPurchase --------\n")
    server := customserver.NewTrainTicketServer()

    request := &pb.PurchaseRequest{
        From:          "StationA",
        To:            "StationB",
        UserFirstName: "John",
        UserLastName:  "Doe",
        UserEmail:     "john.doe@example.com",
    }

    response, err := server.SubmitPurchase(context.Background(), request)

    // Print or log the values
    fmt.Printf("Actual response after SubmitPurchase: %v\n", response)

    assert.NoError(t, err, "Error should be nil")
    assert.NotNil(t, response, "Response should not be nil")
    assert.Equal(t, request.From, response.From, "From should match")
    assert.Equal(t, request.To, response.To, "To should match")
    assert.Equal(t, request.UserFirstName, response.UserFirstName, "UserFirstName should match")
    assert.Equal(t, request.UserLastName, response.UserLastName, "UserLastName should match")
    assert.Equal(t, request.UserEmail, response.UserEmail, "UserEmail should match")
    assert.True(t, response.PricePaid > 0, "PricePaid should be greater than 0")
}

func TestViewReceipt(t *testing.T) {

	
	fmt.Printf("------ TestViewReceipt --------\n")
    server := customserver.NewTrainTicketServer()

    userEmail := "john.doe@example.com"
    server.GetTickets()[userEmail] = &pb.Receipt{
        UserEmail: userEmail,
        // ... other receipt fields ...
    }

    request := &pb.ViewReceiptRequest{
        UserEmail: userEmail,
    }

    response, err := server.ViewReceipt(context.Background(), request)

    // Print or log the values
    fmt.Printf("Actual response after ViewReceipt: %v\n", response)

    assert.NoError(t, err, "Error should be nil")
    assert.NotNil(t, response, "Response should not be nil")
    assert.Equal(t, userEmail, response.UserEmail, "UserEmail should match")
    // ... other assertions for receipt fields ...
}

func TestViewUsersBySection(t *testing.T) {

	fmt.Printf("------ TestViewUsersBySection --------\n")
    server := customserver.NewTrainTicketServer()

    userEmail := "john.doe@example.com"
    section := "A"
    server.GetTickets()[userEmail] = &pb.Receipt{
        UserEmail:   userEmail,
        SeatSection: section,
        // ... other receipt fields ...
    }

    request := &pb.SectionRequest{
        Section: section,
    }

    response, err := server.ViewUsersBySection(context.Background(), request)

    // Print or log the values
    fmt.Printf("Actual response after ViewUsersBySection: %v\n", response)

    assert.NoError(t, err, "Error should be nil")
    assert.NotNil(t, response, "Response should not be nil")
    assert.Len(t, response.UserSeats, 1, "UserSeats should have length 1")
    assert.Equal(t, userEmail, response.UserSeats[0].UserEmail, "UserEmail should match")
    assert.Equal(t, section, response.UserSeats[0].SeatSection, "SeatSection should match")
    // ... other assertions for receipt fields ...
}

func TestRemoveUser(t *testing.T) {

	fmt.Printf("------ TestRemoveUser --------\n")
    server := customserver.NewTrainTicketServer()



    userEmail := "john.doe@example.com"
    section := "A"
    server.GetTickets()[userEmail] = &pb.Receipt{
        UserEmail:   userEmail,
        SeatSection: section,
		SeatNumber: 1,
        // ... other receipt fields ...
    }

	server.GetSeatOccupied()[fmt.Sprintf("%s%d", section, 1)]++

	fmt.Printf("Actual seat occupancy before RemoveUser: %v\n", server.GetSeatOccupied())
    fmt.Printf("Actual tickets before RemoveUser: %v\n", server.GetTickets())


    request := &pb.UserRequest{
        UserEmail: userEmail,
    }


    response, err := server.RemoveUser(context.Background(), request)

    // Print or log the values
    fmt.Printf("Actual response after RemoveUser: %v\n", response)
    fmt.Printf("Actual seat occupancy after RemoveUser: %v\n", server.GetSeatOccupied())
    fmt.Printf("Actual tickets after RemoveUser: %v\n", server.GetTickets())

    assert.NoError(t, err, "Error should be nil")
    assert.NotNil(t, response, "Response should not be nil")
    assert.True(t, response.Success, "Success should be true")
    assert.Equal(t, 0, server.GetSeatOccupied()[fmt.Sprintf("%s%d", section, 1)], "Seat should be unoccupied") // assuming SeatNumber is 1
    assert.Nil(t, server.GetTickets()[userEmail], "User should be removed")
}

func TestModifySeat(t *testing.T) {


	fmt.Printf("------ TestModifySeat --------\n")

    server := customserver.NewTrainTicketServer()

    userEmail := "john.doe@example.com"
    section := "A"



    server.GetTickets()[userEmail] = &pb.Receipt{
        UserEmail:   userEmail,
        SeatSection: section,
		SeatNumber: 1,
        // ... other receipt fields ...
    }
	server.GetSeatOccupied()[fmt.Sprintf("%s%d", section, 1)]++

	fmt.Printf("Actual seat occupancy before RemoveUser: %v\n", server.GetSeatOccupied())
    fmt.Printf("Actual tickets before RemoveUser: %v\n", server.GetTickets())
	

    newSection := "B"
    request := &pb.ModifySeatRequest{
        UserEmail:  userEmail,
        NewSection: newSection,
    }

    response, err := server.ModifySeat(context.Background(), request)


    // Print or log the values
    fmt.Printf("Actual seat occupancy after ModifySeat: %v\n", server.GetSeatOccupied())
    fmt.Printf("Actual tickets after ModifySeat: %v\n", server.GetTickets())

    assert.NoError(t, err, "Error should be nil")
    assert.NotNil(t, response, "Response should not be nil")
    assert.True(t, response.Success, "Success should be true")
    assert.Equal(t, 0, server.GetSeatOccupied()[fmt.Sprintf("%s%d", section, 1)], "Seat should be unoccupied") // assuming SeatNumber is 1
    assert.Equal(t, 1, server.GetSeatOccupied()[fmt.Sprintf("%s%d", newSection, 1)], "New seat should be occupied")
    assert.Equal(t, newSection, server.GetTickets()[userEmail].SeatSection, "SeatSection should be modified")
}

func TestIsSeatOccupied(t *testing.T) {

	fmt.Printf("------ TestIsSeatOccupied --------\n")
    server := customserver.NewTrainTicketServer()

    seatSection := "A1"
    server.GetSeatOccupied()[seatSection] = 1

    occupied := server.IsSeatOccupied(seatSection)
    assert.True(t, occupied, "Seat should be occupied")

    emptySeatSection := "B1"
    occupied = server.IsSeatOccupied(emptySeatSection)
    assert.False(t, occupied, "Seat should be unoccupied")
}
