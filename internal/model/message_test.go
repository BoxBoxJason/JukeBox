package db_model

import "testing"

func TestCreateMessage(t *testing.T) {

	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection to the database: %v", err)
	}
	defer CloseConnection(db)

	// Create a user for the message
	user := &User{
		Email:           "test_user_200@gmail.com",
		Hashed_Password: "hashed_password",
		Username:        "test_user_200",
	}

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Create a new message
	message := &Message{
		Content: "test_content_1",
		Sender:  *user,
	}

	err = message.CreateMessage(db)
	if err != nil {
		t.Errorf("Error creating message: %v", err)
	}
}

func TestCreateMessages(t *testing.T) {

	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection to the database: %v", err)
	}
	defer CloseConnection(db)

	// Create a user for the messages
	user := &User{
		Email:           "test_user_201@gmail.com",
		Hashed_Password: "hashed_password",
		Username:        "test_user_201",
	}

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Create multiple messages
	messages := []*Message{
		{
			Content: "test_content_2",
			Sender:  *user,
		},
		{
			Content: "test_content_3",
			Sender:  *user,
		},
	}

	err = CreateMessages(db, messages)
	if err != nil {
		t.Errorf("Error creating messages: %v", err)
	}
}

func TestGetMessageByID(t *testing.T) {

	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection to the database: %v", err)
	}
	defer CloseConnection(db)

	// Create a user for the message
	user := &User{
		Email:           "test_user_202@gmail.com",
		Hashed_Password: "hashed_password",
		Username:        "test_user_202",
	}

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Create a new message
	message := &Message{
		Content: "test_content_4",
		Sender:  *user,
	}

	err = message.CreateMessage(db)
	if err != nil {
		t.Errorf("Error creating message: %v", err)
	}

	// Retrieve the message by ID
	_, err = GetMessageByID(db, message.ID)
	if err != nil {
		t.Errorf("Error retrieving message by ID: %v", err)
	}
}

func TestGetMessagesBySender(t *testing.T) {

	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection to the database: %v", err)
	}
	defer CloseConnection(db)

	// Create a user for the messages
	user := &User{
		Email:           "test_user_203@gmail.com",
		Hashed_Password: "hashed_password",
		Username:        "test_user_203",
	}

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Create multiple messages
	messages := []*Message{
		{
			Content: "test_content_5",
			Sender:  *user,
		},
		{
			Content: "test_content_6",
			Sender:  *user,
		},
	}

	err = CreateMessages(db, messages)
	if err != nil {
		t.Errorf("Error creating messages: %v", err)
	}

	// Retrieve the messages by sender
	retrieved_messages, err := user.GetMessages(db)
	if err != nil {
		t.Errorf("Error retrieving messages by sender: %v", err)
	}
	if len(retrieved_messages) < 2 {
		t.Errorf("Error retrieving messages by sender: expected at least 2 messages, got %v", len(retrieved_messages))
	}
}

func TestGetAllMessagesBySender(t *testing.T) {

	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection to the database: %v", err)
	}
	defer CloseConnection(db)

	// Create a user for the messages
	user := &User{
		Email:           "test_user_204@gmail.com",
		Hashed_Password: "hashed_password",
		Username:        "test_user_204",
	}

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Create multiple messages
	messages := []*Message{
		{
			Content:  "test_content_7",
			Sender:   *user,
			Censored: true,
		},
		{
			Content:  "test_content_8",
			Sender:   *user,
			Censored: true,
		},
	}

	err = CreateMessages(db, messages)
	if err != nil {
		t.Errorf("Error creating messages: %v", err)
	}

	// Retrieve the messages by sender
	retrieved_messages, err := user.GetAllMessages(db)
	if err != nil {
		t.Errorf("Error retrieving messages by sender: %v", err)
	}
	if len(retrieved_messages) < 2 {
		t.Errorf("Error retrieving messages by sender: expected at least 2 messages, got %v", len(retrieved_messages))
	}
}

func TestGetMessagesByContent(t *testing.T) {

	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection to the database: %v", err)
	}
	defer CloseConnection(db)

	// Create a user for the messages
	user := &User{
		Email:           "test_user_205@gmail.com",
		Hashed_Password: "hashed_password",
		Username:        "test_user_205",
	}

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Create multiple messages
	messages := []*Message{
		{
			Content: "test_content_9",
			Sender:  *user,
		},
		{
			Content: "test_content_10",
			Sender:  *user,
		},
	}

	err = CreateMessages(db, messages)
	if err != nil {
		t.Errorf("Error creating messages: %v", err)
	}

	// Retrieve the messages by content
	retrieved_messages, err := GetMessagesByContent(db, "test_content")
	if err != nil {
		t.Errorf("Error retrieving messages by content: %v", err)
	}
	if len(retrieved_messages) < 2 {
		t.Errorf("Error retrieving messages by content: expected at least 2 message, got %v", len(retrieved_messages))
	}
}

func TestGetAllMessagesByContent(t *testing.T) {

	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection to the database: %v", err)
	}
	defer CloseConnection(db)

	// Create a user for the messages
	user := &User{
		Email:           "test_user_206@gmail.com",
		Hashed_Password: "hashed_password",
		Username:        "test_user_206",
	}

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Create multiple messages
	messages := []*Message{
		{
			Content:  "test_content_11",
			Sender:   *user,
			Censored: true,
		},
		{
			Content:  "test_content_12",
			Sender:   *user,
			Censored: true,
		},
	}

	err = CreateMessages(db, messages)
	if err != nil {
		t.Errorf("Error creating messages: %v", err)
	}

	// Retrieve the messages by content
	retrieved_messages, err := GetAllMessagesByContent(db, "test_content")
	if err != nil {
		t.Errorf("Error retrieving messages by content: %v", err)
	}
	if len(retrieved_messages) < 2 {
		t.Errorf("Error retrieving messages by content: expected at least 2 message, got %v", len(retrieved_messages))
	}
}

func TestGetAllVisibleMessages(t *testing.T) {

	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection to the database: %v", err)
	}
	defer CloseConnection(db)

	// Create a user for the messages
	user := &User{
		Email:           "test_user_207@gmail.com",
		Hashed_Password: "hashed_password",
		Username:        "test_user_207",
	}

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Create multiple messages
	messages := []*Message{
		{
			Content: "test_content_13",
			Sender:  *user,
		},
		{
			Content: "test_content_14",
			Sender:  *user,
		},
	}

	err = CreateMessages(db, messages)
	if err != nil {
		t.Errorf("Error creating messages: %v", err)
	}

	// Retrieve all visible messages
	retrieved_messages, err := GetAllVisibleMessages(db)
	if err != nil {
		t.Errorf("Error retrieving all visible messages: %v", err)
	}
	if len(retrieved_messages) < 2 {
		t.Errorf("Error retrieving all visible messages: expected at least 2 messages, got %v", len(retrieved_messages))
	}
}

func TestGetAllMessages(t *testing.T) {

	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection to the database: %v", err)
	}
	defer CloseConnection(db)

	// Create a user for the messages
	user := &User{
		Email:           "test_user_208@gmail.com",
		Hashed_Password: "hashed_password",
		Username:        "test_user_208",
	}

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Create multiple messages
	messages := []*Message{
		{
			Content: "test_content_15",
			Sender:  *user,
		},
		{
			Content: "test_content_16",
			Sender:  *user,
		},
	}

	err = CreateMessages(db, messages)
	if err != nil {
		t.Errorf("Error creating messages: %v", err)
	}

	// Retrieve all messages
	retrieved_messages, err := GetAllMessages(db)
	if err != nil {
		t.Errorf("Error retrieving all messages: %v", err)
	}
	if len(retrieved_messages) < 2 {
		t.Errorf("Error retrieving all messages: expected at least 2 messages, got %v", len(retrieved_messages))
	}
}

func TestGetFlaggedMessages(t *testing.T) {

	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection to the database: %v", err)
	}
	defer CloseConnection(db)

	// Create a user for the messages
	user := &User{
		Email:           "test_user_209@gmail.com",
		Hashed_Password: "hashed_password",
		Username:        "test_user_209",
	}

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Create multiple messages
	messages := []*Message{
		{
			Content: "test_content_17",
			Sender:  *user,
			Flagged: true,
		},
		{
			Content: "test_content_18",
			Sender:  *user,
			Flagged: true,
		},
	}

	err = CreateMessages(db, messages)
	if err != nil {
		t.Errorf("Error creating messages: %v", err)
	}

	// Retrieve the flagged messages
	retrieved_messages, err := GetFlaggedMessages(db)
	if err != nil {
		t.Errorf("Error retrieving flagged messages: %v", err)
	}
	if len(retrieved_messages) < 2 {
		t.Errorf("Error retrieving flagged messages: expected at least 2 messages, got %v", len(retrieved_messages))
	}
}

func TestGetFlaggedMessagesBySender(t *testing.T) {

	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection to the database: %v", err)
	}
	defer CloseConnection(db)

	// Create a user for the messages
	user := &User{
		Email:           "test_user_210@gmail.com",
		Hashed_Password: "hashed_password",
		Username:        "test_user_210",
	}

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Create multiple messages
	messages := []*Message{
		{
			Content: "test_content_19",
			Sender:  *user,
			Flagged: true,
		},
		{
			Content: "test_content_20",
			Sender:  *user,
			Flagged: true,
		},
	}

	err = CreateMessages(db, messages)
	if err != nil {
		t.Errorf("Error creating messages: %v", err)
	}

	// Retrieve the flagged messages by sender
	retrieved_messages, err := user.GetFlaggedMessages(db)
	if err != nil {
		t.Errorf("Error retrieving flagged messages by sender: %v", err)
	}
	if len(retrieved_messages) < 2 {
		t.Errorf("Error retrieving flagged messages by sender: expected at least 2 messages, got %v", len(retrieved_messages))
	}
}

func TestGetRemovedMessages(t *testing.T) {

	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection to the database: %v", err)
	}
	defer CloseConnection(db)

	// Create a user for the messages
	user := &User{
		Email:           "test_user_211@gmail.com",
		Hashed_Password: "hashed_password",
		Username:        "test_user_211",
	}

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Create multiple messages
	messages := []*Message{
		{
			Content: "test_content_21",
			Sender:  *user,
			Removed: true,
		},
		{
			Content: "test_content_22",
			Sender:  *user,
			Removed: true,
		},
	}

	err = CreateMessages(db, messages)
	if err != nil {
		t.Errorf("Error creating messages: %v", err)
	}

	// Retrieve the removed messages
	retrieved_messages, err := GetRemovedMessages(db)
	if err != nil {
		t.Errorf("Error retrieving removed messages: %v", err)
	}
	if len(retrieved_messages) < 2 {
		t.Errorf("Error retrieving removed messages: expected at least 2 messages, got %v", len(retrieved_messages))
	}
}

func TestGetRemovedMessagesBySender(t *testing.T) {

	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection to the database: %v", err)
	}
	defer CloseConnection(db)

	// Create a user for the messages
	user := &User{
		Email:           "test_user_212@gmail.com",
		Hashed_Password: "hashed_password",
		Username:        "test_user_212",
	}

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Create multiple messages
	messages := []*Message{
		{
			Content: "test_content_23",
			Sender:  *user,
			Removed: true,
		},
		{
			Content: "test_content_24",
			Sender:  *user,
			Removed: true,
		},
	}

	err = CreateMessages(db, messages)
	if err != nil {
		t.Errorf("Error creating messages: %v", err)
	}

	// Retrieve the removed messages by sender
	retrieved_messages, err := user.GetRemovedMessages(db)
	if err != nil {
		t.Errorf("Error retrieving removed messages by sender: %v", err)
	}
	if len(retrieved_messages) < 2 {
		t.Errorf("Error retrieving removed messages by sender: expected at least 2 messages, got %v", len(retrieved_messages))
	}
}

func TestGetCensoredMessages(t *testing.T) {

	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection to the database: %v", err)
	}
	defer CloseConnection(db)

	// Create a user for the messages
	user := &User{
		Email:           "test_user_213@gmail.com",
		Hashed_Password: "hashed_password",
		Username:        "test_user_213",
	}

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Create multiple messages
	messages := []*Message{
		{
			Content:  "test_content_25",
			Sender:   *user,
			Censored: true,
		},
		{
			Content:  "test_content_26",
			Sender:   *user,
			Censored: true,
		},
	}

	err = CreateMessages(db, messages)
	if err != nil {
		t.Errorf("Error creating messages: %v", err)
	}

	// Retrieve the censored messages
	retrieved_messages, err := GetCensoredMessages(db)
	if err != nil {
		t.Errorf("Error retrieving censored messages: %v", err)
	}
	if len(retrieved_messages) < 2 {
		t.Errorf("Error retrieving censored messages: expected at least 2 messages, got %v", len(retrieved_messages))
	}
}

func TestGetCensoredMessagesBySender(t *testing.T) {
	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection to the database: %v", err)
	}
	defer CloseConnection(db)

	// Create a user for the messages
	user := &User{
		Email:           "test_user_214@gmail.com",
		Hashed_Password: "hashed_password",
		Username:        "test_user_214",
	}

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Create multiple messages
	messages := []*Message{
		{
			Content:  "test_content_27",
			Sender:   *user,
			Censored: true,
		},
		{
			Content:  "test_content_28",
			Sender:   *user,
			Censored: true,
		},
	}

	err = CreateMessages(db, messages)
	if err != nil {
		t.Errorf("Error creating messages: %v", err)
	}

	// Retrieve the censored messages by sender
	retrieved_messages, err := user.GetCensoredMessages(db)
	if err != nil {
		t.Errorf("Error retrieving censored messages by sender: %v", err)
	}
	if len(retrieved_messages) < 2 {
		t.Errorf("Error retrieving censored messages by sender: expected at least 2 messages, got %v", len(retrieved_messages))
	}
}

func TestUpdateMessage(t *testing.T) {
	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection to the database: %v", err)
	}
	defer CloseConnection(db)

	// Create a user for the message
	user := &User{
		Email:           "test_user_215@gmail.com",
		Hashed_Password: "hashed_password",
		Username:        "test_user_215",
	}

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Create a new message
	message := &Message{
		Content: "test_content_29",
		Sender:  *user,
	}

	err = message.CreateMessage(db)
	if err != nil {
		t.Errorf("Error creating message: %v", err)
	}

	// Update the message
	message.Content = "test_content_29_updated"

	err = message.UpdateMessage(db)
	if err != nil {
		t.Errorf("Error updating message: %v", err)
	}

	// Retrieve the updated message
	retrieved_message, err := GetMessageByID(db, message.ID)
	if err != nil {
		t.Errorf("Error retrieving message by ID: %v", err)
	}

	if retrieved_message.Content != message.Content {
		t.Errorf("Updated message does not match")
	}
}

func TestDeleteMessage(t *testing.T) {
	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection to the database: %v", err)
	}
	defer CloseConnection(db)

	// Create a user for the message
	user := &User{
		Email:           "test_user_216@gmail.com",
		Hashed_Password: "hashed_password",
		Username:        "test_user_216",
	}

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Create a new message
	message := &Message{
		Content: "test_content_30",
		Sender:  *user,
	}

	err = message.CreateMessage(db)
	if err != nil {
		t.Errorf("Error creating message: %v", err)
	}

	// Delete the message
	err = message.DeleteMessage(db)
	if err != nil {
		t.Errorf("Error deleting message: %v", err)
	}

	// Retrieve the deleted message
	_, err = GetMessageByID(db, message.ID)
	if err == nil {
		t.Errorf("Error retrieving deleted message by ID: %v", err)
	}
}

func TestDeleteMessages(t *testing.T) {
	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection to the database: %v", err)
	}
	defer CloseConnection(db)

	// Create a user for the messages
	user := &User{
		Email:           "test_user_217@gmail.com",
		Hashed_Password: "hashed_password",
		Username:        "test_user_217",
	}

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Create multiple messages
	messages := []*Message{
		{
			Content: "test_content_31",
			Sender:  *user,
		},
		{
			Content: "test_content_32",
			Sender:  *user,
		},
	}

	err = CreateMessages(db, messages)
	if err != nil {
		t.Errorf("Error creating messages: %v", err)
	}

	// Delete the messages
	err = DeleteMessages(db, messages)
	if err != nil {
		t.Errorf("Error deleting messages: %v", err)
	}

	// Retrieve the deleted messages
	for _, message := range messages {
		_, err = GetMessageByID(db, message.ID)
		if err == nil {
			t.Errorf("Error retrieving deleted message by ID: %v", err)
		}
	}
}
