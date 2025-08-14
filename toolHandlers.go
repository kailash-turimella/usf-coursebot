package main

import (
	"encoding/json"
	"fmt"
)

type toolHandler func(collection CollectionData, args string) string

var handlerMap = map[string]toolHandler{
	"get_courses_by_instructor":       getCoursesByInstructorToolHandler,
	"get_courses_by_subject":          getCoursesBySubjectToolHandler,
	"get_courses_by_title":            getCoursesByTitleToolHandler,
	"get_courses_by_meeting_day":      getCoursesByMeetingDayToolHandler,
	"get_courses_by_instruction_mode": getCoursesByInstructionModeToolHandler,
	"get_courses_by_time":             getCoursesByTimeToolHandler,
	"get_calendar_invite":             getCalendarInviteToolHandler,
}

func getCoursesByInstructorToolHandler(collection CollectionData, args string) string {
	var parsed struct {
		InstructorName string `json:"instructor_name"`
	}
	if err := json.Unmarshal([]byte(args), &parsed); err != nil {
		fmt.Printf("Error reading the json")
		return ""
	}

	if !useTerminal {
		fmt.Printf("Looking for Instructor Name: %s\n", parsed.InstructorName)
	}

	where := map[string]any{"Primary Instructor Last Name": parsed.InstructorName}

	results, err := collection.collection.Query(collection.context, []string{"courses taught by " + parsed.InstructorName}, numRelevantCoursesToRetrieve, where, nil, nil)
	if err != nil {
		fmt.Printf("Error querying the database: %v\n", err)
		return ""
	}

	return resultsToString(results)
}

func getCoursesBySubjectToolHandler(collection CollectionData, args string) string {
	var parsed struct {
		Subject string `json:"subject"`
	}
	if err := json.Unmarshal([]byte(args), &parsed); err != nil {
		fmt.Printf("Error reading the json")
		return ""
	}

	if !useTerminal {
		fmt.Printf("Looking for Subject: %s\n", parsed.Subject)
	}

	where := map[string]any{"SUBJ": parsed.Subject}

	results, err := collection.collection.Query(collection.context, []string{"courses for subject " + parsed.Subject}, numRelevantCoursesToRetrieve, where, nil, nil)
	if err != nil {
		fmt.Printf("Error querying the database: %v\n", err)
		return ""
	}

	return resultsToString(results)
}

func getCoursesByTitleToolHandler(collection CollectionData, args string) string {
	var parsed struct {
		TitleKeyword string `json:"title_keyword"`
	}
	if err := json.Unmarshal([]byte(args), &parsed); err != nil {
		fmt.Printf("Error reading the json")
		return ""
	}

	if !useTerminal {
		fmt.Printf("Looking for Title Keyword: %s\n", parsed.TitleKeyword)
	}

	where := map[string]any{"Title Short Desc": parsed.TitleKeyword}

	results, err := collection.collection.Query(collection.context, []string{"courses with title " + parsed.TitleKeyword}, numRelevantCoursesToRetrieve, where, nil, nil)
	if err != nil {
		fmt.Printf("Error querying the database: %v\n", err)
		return ""
	}

	return resultsToString(results)
}

func getCoursesByMeetingDayToolHandler(collection CollectionData, args string) string {
	var parsed struct {
		MeetDays string `json:"meet_days"`
	}
	if err := json.Unmarshal([]byte(args), &parsed); err != nil {
		fmt.Printf("Error reading the json\n")
		return ""
	}

	if !useTerminal {
		fmt.Printf("Looking for Meeting Days: %s\n", parsed.MeetDays)
	}

	where := map[string]any{"Meet Days": parsed.MeetDays}

	results, err := collection.collection.Query(collection.context, []string{"courses on days " + parsed.MeetDays}, numRelevantCoursesToRetrieve, where, nil, nil)
	if err != nil {
		fmt.Printf("Error querying the database: %v\n", err)
		return ""
	}

	return resultsToString(results)
}

func getCoursesByInstructionModeToolHandler(collection CollectionData, args string) string {
	var parsed struct {
		InstructionMode string `json:"instruction_mode"`
	}
	if err := json.Unmarshal([]byte(args), &parsed); err != nil {
		fmt.Printf("Error reading the json\n")
		return ""
	}

	if !useTerminal {
		fmt.Printf("Looking for Instruction Mode: %s\n", parsed.InstructionMode)
	}

	where := map[string]any{"Instruction Mode Desc": parsed.InstructionMode}

	results, err := collection.collection.Query(collection.context, []string{"courses with instruction mode " + parsed.InstructionMode}, numRelevantCoursesToRetrieve, where, nil, nil)
	if err != nil {
		fmt.Printf("Error querying the database: %v\n", err)
		return ""
	}

	return resultsToString(results)
}

func getCoursesByTimeToolHandler(collection CollectionData, args string) string {
	var parsed struct {
		BeforeOrAfter string `json:"before_or_after"`
		Time          string `json:"time"`
	}
	if err := json.Unmarshal([]byte(args), &parsed); err != nil {
		fmt.Printf("Error reading the json\n")
		return ""
	}

	if !useTerminal {
		fmt.Printf("Looking for Time: %s %s\n", parsed.BeforeOrAfter, parsed.Time)
	}

	// Filtering manually since Chroma doesn't support inequality in metadata
	results, err := collection.collection.Query(collection.context, []string{"courses filtered by time"}, numRelevantCoursesToRetrieve*1000, nil, nil, nil)
	if err != nil {
		fmt.Printf("Error querying the database: %v\n", err)
		return ""
	}

	filtered := []int{}
	for i, meta := range results.Metadatas[0] {
		timeStr, ok := meta["Begin Time"].(string)
		if !ok || len(timeStr) < 4 {
			continue
		}
		if parsed.BeforeOrAfter == "before" && timeStr < parsed.Time {
			filtered = append(filtered, i)
		} else if parsed.BeforeOrAfter == "after" && timeStr > parsed.Time {
			filtered = append(filtered, i)
		}
	}

	courses := ""
	for _, i := range filtered {
		courses += "Document: " + results.Documents[0][i] + "\n"
		courses += "Metadata:\n" + getMetadata(results.Metadatas[0][i]) + "\n"
	}
	return courses
}

func getCalendarInviteToolHandler(collection CollectionData, args string) string {
	var parsed struct {
		Course string `json:"course"`
		Format string `json:"format"`
	}
	if err := json.Unmarshal([]byte(args), &parsed); err != nil {
		fmt.Printf("Error reading the json\n")
		return ""
	}

	if !useTerminal {
		fmt.Printf("Generating calendar invite for Course: %s in format: %s\n", parsed.Course, parsed.Format)
	}

	// Remove the where clause to search across all courses
	results, err := collection.collection.Query(collection.context, []string{parsed.Course}, 4, nil, nil, nil)
	if err != nil {
		fmt.Printf("Error querying the database: %v\n", err)
		return ""
	}

	if len(results.Documents) == 0 || len(results.Documents[0]) == 0 {
		fmt.Printf("Course not found: %s\n", parsed.Course)
		return "Error: Course not found"
	}

	metadata := results.Metadatas[0][0]

	// Extract course details with error checking
	title, ok := metadata["Title Short Desc"].(string)
	if !ok {
		fmt.Printf("Missing or invalid title in metadata\n")
		return "Error: Missing course title"
	}

	meetDays, ok := metadata["Meet Days"].(string)
	if !ok {
		fmt.Printf("Missing or invalid meeting days in metadata\n")
		return "Error: Missing meeting days"
	}

	beginTime, ok := metadata["Begin Time"].(string)
	if !ok {
		fmt.Printf("Missing or invalid begin time in metadata\n")
		return "Error: Missing begin time"
	}

	endTime, ok := metadata["End Time"].(string)
	if !ok {
		fmt.Printf("Missing or invalid end time in metadata\n")
		return "Error: Missing end time"
	}

	location, ok := metadata["BLDG"].(string)
	location += " " + metadata["RM"].(string)
	if !ok {
		fmt.Printf("Missing or invalid location in metadata\n")
		return "Error: Missing location"
	}

	instructor, ok := metadata["Primary Instructor Last Name"].(string)
	if !ok {
		fmt.Printf("Missing or invalid instructor in metadata\n")
		return "Error: Missing instructor"
	}

	// Format time from "HHMM" to "HH:MM"
	formatTime := func(timeStr string) string {
		if len(timeStr) < 4 {
			return timeStr
		}
		return timeStr[:2] + ":" + timeStr[2:]
	}

	if parsed.Format == "details" {
		return fmt.Sprintf("Course: %s\nInstructor: %s\nDays: %s\nTime: %s - %s\nLocation: %s",
			title, instructor, meetDays, formatTime(beginTime), formatTime(endTime), location)
	}

	// Generate .ics format
	// Note: This is a simplified .ics format. A full implementation would need proper date handling
	// and more complete .ics specification compliance
	ics := fmt.Sprintf(`BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
SUMMARY:%s
DESCRIPTION:Taught by %s
LOCATION:%s
DTSTART;TZID=America/Los_Angeles:20240101T%s00
DTEND;TZID=America/Los_Angeles:20240101T%s00
RRULE:FREQ=WEEKLY;BYDAY=%s
END:VEVENT
END:VCALENDAR`,
		title, instructor, location, formatTime(beginTime), formatTime(endTime), meetDays)

	return ics
}
