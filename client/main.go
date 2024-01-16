package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"strings"

	"google.golang.org/grpc"
	pb "github.com/pkkp86nitd/train_ticket_app/proto" 
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTrainTicketServiceClient(conn)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Choose an action:")
		fmt.Println("1. Submit Purchase")
		fmt.Println("2. View Receipt")
		fmt.Println("3. View Users By Section")
		fmt.Println("4. Remove User")
		fmt.Println("5. Modify Seat")
		fmt.Println("0. Exit")

		fmt.Print("Enter your choice (0-5): ")
		choice, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("error reading input: %v", err)
		}

		switch choice {
		case "1\n":
			submitPurchase(client, reader)
		case "2\n":
			viewReceipt(client, reader)
		case "3\n":
			viewUsersBySection(client, reader)
		case "4\n":
			removeUser(client, reader)
		case "5\n":
			modifySeat(client, reader)
		case "0\n":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice. Please enter a valid option.")
		}
	}
}

func submitPurchase(client pb.TrainTicketServiceClient, reader *bufio.Reader) {
	fmt.Println("Submitting Purchase:")
	from := getInput("From: ", reader)
	to := getInput("To: ", reader)
	firstName := getInput("User First Name: ", reader)
	lastName := getInput("User Last Name: ", reader)
	email := getInput("User Email: ", reader)

	response, err := client.SubmitPurchase(context.TODO(), &pb.PurchaseRequest{
		From:          from,
		To:            to,
		UserFirstName: firstName,
		UserLastName:  lastName,
		UserEmail:     email,
	})
	if err != nil {
		fmt.Printf("could not submit purchase: %v\n", err)
		return
	}
	fmt.Printf("Purchase submitted successfully. Receipt: %+v\n", response)
	time.Sleep(1 * time.Second) // Add a delay of 1 second
}

func viewReceipt(client pb.TrainTicketServiceClient, reader *bufio.Reader) {
	fmt.Println("Viewing Receipt:")
	email := getInput("User Email: ", reader)

	response, err := client.ViewReceipt(context.TODO(), &pb.ViewReceiptRequest{
		UserEmail: email,
	})
	if err != nil {
		fmt.Printf("could not view receipt: %v\n", err)
		return
	}
	fmt.Printf("View receipt response: %+v\n", response)
	time.Sleep(1 * time.Second) // Add a delay of 1 second
}

func viewUsersBySection(client pb.TrainTicketServiceClient, reader *bufio.Reader) {
	fmt.Println("Viewing Users By Section:")
	section := getInput("Section: ", reader)

	response, err := client.ViewUsersBySection(context.TODO(), &pb.SectionRequest{
		Section: section,
	})
	if err != nil {
		fmt.Printf("could not view users by section: %v\n", err)
		return
	}
	fmt.Printf("View users by section response: %+v\n", response)
	time.Sleep(1 * time.Second) // Add a delay of 1 second
}

func removeUser(client pb.TrainTicketServiceClient, reader *bufio.Reader) {
	fmt.Println("Removing User:")
	email := getInput("User Email: ", reader)

	response, err := client.RemoveUser(context.TODO(), &pb.UserRequest{
		UserEmail: email,
	})
	if err != nil {
		fmt.Printf("could not remove user: %v\n", err)
		return
	}
	fmt.Printf("Remove user response: %+v\n", response)
	time.Sleep(1 * time.Second) // Add a delay of 1 second
}

func modifySeat(client pb.TrainTicketServiceClient, reader *bufio.Reader) {
	fmt.Println("Modifying Seat:")
	email := getInput("User Email: ", reader)
	newSection := getInput("New Section: ", reader)

	response, err := client.ModifySeat(context.TODO(), &pb.ModifySeatRequest{
		UserEmail:  email,
		NewSection: newSection,
	})
	if err != nil {
		fmt.Printf("could not modify seat: %v\n", err)
		return
	}
	fmt.Printf("Modify seat response: %+v\n", response)
	time.Sleep(1 * time.Second) // Add a delay of 1 second
}

func getInput(prompt string, reader *bufio.Reader) string {
	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("error reading input: %v", err)
	}
	return strings.TrimSpace(input)
}
