package assistant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gitto/features"
	"io"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

const OLLAMA_API_URL = "http://localhost:11434/api/generate"

type LLMAssistant struct {
}

type InterpretResponse struct {
	Action       string `json:"action"`
	ReminderText string `json:"reminder_text"`
	ReminderDate string `json:"reminder_date"`
	NoteTitle    string `json:"note_title"`
	NoteContent  string `json:"note_content"`
	EventName    string `json:"event_name"`
	EventStart   string `json:"event_start"`
	EventEnd     string `json:"event_end"`
	TodoTask     string `json:"todo_task"`
	TodoID       string `json:"todo_id"`
	UpdateField  string `json:"update_field"`
	UpdateValue  string `json:"update_value"`
}

func (a *LLMAssistant) sendRequest(prompt string) (string, error) {
	data := map[string]interface{}{
		"model":  "llama3",
		"prompt": prompt,
		"stream": true,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", OLLAMA_API_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	println("Send")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	println("Recv")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	var message string
	decoder := json.NewDecoder(resp.Body)

	for {
		var body map[string]interface{}
		if err := decoder.Decode(&body); err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}

		if errMsg, exists := body["error"]; exists {
			return "", fmt.Errorf("%v", errMsg)
		}

		if done, exists := body["done"]; exists && done == true {
			return message, nil
		}

		if response, exists := body["response"].(string); exists {
			message += response
		}
	}

	return message, nil
}

func (a *LLMAssistant) queryLLM(prompt string) (string, error) {
	return a.sendRequest(prompt)
}

func (a *LLMAssistant) interpretInput(userInput string) (*InterpretResponse, error) {
	prompt := fmt.Sprintf(`
	Given the user information and the following input, interpret the user's intent and categorize it as one of these actions: REMIND, NOTE, SCHEDULE, TODO, VIEW_TODOS, COMPLETE_TODO, DELETE_TODO, UPDATE_INFO, SUGGEST_PRIORITIES, or UNKNOWN.
	For REMIND, extract the reminder text and the exact time expression used (e.g., "in 5 minutes", "tomorrow at 3pm").
	For NOTE, extract the title and content.
	For SCHEDULE, extract the event name, start time, and end time. Provide the exact time expressions used.
	For TODO, extract the task.
	For COMPLETE_TODO or DELETE_TODO, extract the todo ID.
	For UPDATE_INFO, identify which user information field needs to be updated.
	For SUGGEST_PRIORITIES, no additional information is needed.

	User input: "%s"

	Your response must be in valid JSON format as shown below. Do not include any other text outside the JSON structure.
	{
		"action": "REMIND/NOTE/VIEW_NOTES/DELETE_NOTE/SCHEDULE/TODO/VIEW_TODOS/COMPLETE_TODO/DELETE_TODO/UPDATE_INFO/SUGGEST_PRIORITIES/UNKNOWN",
		"reminder_text": "",
		"reminder_date": "",
		"note_title": "",
		"note_content": "",
		"event_name": "",
		"event_start": "",
		"event_end": "",
		"todo_task": "",
		"todo_id": "",
		"update_field": "",
		"update_value": ""
	}
	`, userInput)

	response, err := a.queryLLM(prompt)
	if err != nil {
		return nil, err
	}
	println(response)
	var interpreted InterpretResponse
	if err := json.Unmarshal([]byte(response), &interpreted); err != nil {
		fmt.Println("Error: Unable to parse LLM response. Using fallback interpretation.")
		return &InterpretResponse{Action: "UNKNOWN"}, nil
	}

	// Handle date parsing as needed, replace with appropriate Go date handling
	// e.g., parsedDate, err := time.Parse(time.RFC3339, interpreted.ReminderDate)

	return &interpreted, nil
}

func (a *LLMAssistant) Process(db *gorm.DB, userInput string) {
	println("hey")
	interpreted, e := a.interpretInput(userInput)
	println("bey")
	if e != nil {
		fmt.Println("Error:", e)
		return
	}

	var err error

	switch interpreted.Action {
	case "REMIND":
		// Implement reminders.add_reminder logic
	case "NOTE":
		// Implement notes.add_note logic
	case "VIEW_NOTES":
		// Implement notes.view_notes logic
	case "DELETE_NOTE":
		// Implement notes.delete_note logic
	case "SCHEDULE":
		// Implement schedule.add_schedule logic
	case "TODO":
		err = features.SaveTodo(db, interpreted.TodoTask)
	case "VIEW_TODOS":
		todos, err := features.ViewTodos(db)
		if err == nil {
			fmt.Println("Your todos:")
			for _, todo := range todos {
				fmt.Printf("Id: %d, Task: %s, Status: %s\n", todo.ID, todo.Task, todo.Status)
			}
		}
	case "COMPLETE_TODO":
		id := interpreted.TodoID
		num_id, err := strconv.Atoi(id)
		err = features.UpdateTodo(db, uint(num_id))
		if err == nil {
			todos, err := features.ViewTodos(db)

			if err == nil {
				fmt.Println("Your todos:")
				for _, todo := range todos {
					fmt.Printf("Id: %d, Task: %s, Status: %s\n", todo.ID, todo.Task, todo.Status)
				}
			}
		}

	case "DELETE_TODO":
		// Implement todos.delete_todo logic
	case "UPDATE_INFO":
		// Implement user_info.update_user_info logic
	case "SUGGEST_PRIORITIES":
		// Implement todos.suggest_priority_tasks logic
		fmt.Println("Here are my suggestions for task prioritization:")
		// Print suggestions
	default:
		fmt.Println("I'm not sure what you want me to do. Can you please rephrase your request?")
	}

	if err != nil {
		fmt.Println("Errror:", err)
	}
}
