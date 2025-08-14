package main

import (
	"github.com/sashabaranov/go-openai"
)

func getAllTools() []openai.Tool {
	return []openai.Tool{
		getCoursesByInstructorTool,
		getCoursesBySubjectTool,
		getCoursesByTitleTool,
		getCoursesByMeetingDayTool,
		getCoursesByInstructionModeTool,
		getCoursesByTimeTool,
		getCalendarInviteTool,
	}
}

var getCoursesByInstructorTool = openai.Tool{
	Type: openai.ToolTypeFunction,
	Function: &openai.FunctionDefinition{
		Name:        "get_courses_by_instructor",
		Description: "Retrieve all courses taught by a specific instructor",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"instructor_name": map[string]any{
					"type":        "string",
					"description": "The last name of the instructor (e.g., 'Peterson')",
				},
			},
			"required": []string{"instructor_name"},
		},
	},
}

var getCoursesBySubjectTool = openai.Tool{
	Type: openai.ToolTypeFunction,
	Function: &openai.FunctionDefinition{
		Name:        "get_courses_by_subject",
		Description: "Retrieve all courses of a given subject abbreviation",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"subject": map[string]any{
					"type":        "string",
					"description": "The subject abbreviation (e.g., 'PHIL', 'CS')",
				},
			},
			"required": []string{"subject"},
		},
	},
}

var getCoursesByTitleTool = openai.Tool{
	Type: openai.ToolTypeFunction,
	Function: &openai.FunctionDefinition{
		Name:        "get_courses_by_title",
		Description: "Retrieve all courses with a given title or keyword in their title",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"title_keyword": map[string]any{
					"type":        "string",
					"description": "A keyword or full title to match in the course title",
				},
			},
			"required": []string{"title_keyword"},
		},
	},
}

var getCoursesByMeetingDayTool = openai.Tool{
	Type: openai.ToolTypeFunction,
	Function: &openai.FunctionDefinition{
		Name:        "get_courses_by_meeting_day",
		Description: "Retrieve all courses that meet on specific days",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"meet_days": map[string]any{
					"type":        "string",
					"description": "The meeting days string (e.g., 'MWF', 'TR', 'MW')",
				},
			},
			"required": []string{"meet_days"},
		},
	},
}

var getCoursesByInstructionModeTool = openai.Tool{
	Type: openai.ToolTypeFunction,
	Function: &openai.FunctionDefinition{
		Name:        "get_courses_by_instruction_mode",
		Description: "Retrieve courses by their instruction mode, e.g., In-Person, Online, or Hybrid",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"instruction_mode": map[string]any{
					"type":        "string",
					"description": "The instruction mode (e.g., 'Online', 'Hybrid', 'In-Person')",
				},
			},
			"required": []string{"instruction_mode"},
		},
	},
}

var getCoursesByTimeTool = openai.Tool{
	Type: openai.ToolTypeFunction,
	Function: &openai.FunctionDefinition{
		Name:        "get_courses_by_time",
		Description: "Retrieve all courses that start before or after a specific time",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"before_or_after": map[string]any{
					"type":        "string",
					"description": "'before' or 'after'",
				},
				"time": map[string]any{
					"type":        "string",
					"description": "Time in 24hr format without colon (e.g., '1200' for 12:00 PM)",
				},
			},
			"required": []string{"before_or_after", "time"},
		},
	},
}

var getCalendarInviteTool = openai.Tool{
	Type: openai.ToolTypeFunction,
	Function: &openai.FunctionDefinition{
		Name:        "get_calendar_invite",
		Description: "Generate a calendar invite (.ics format) or event details for a course",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"course": map[string]any{
					"type":        "string",
					"description": "The course to generate calendar invite for",
				},
				"format": map[string]any{
					"type":        "string",
					"description": "The output format ('ics' for .ics file or 'details' for human-readable details)",
					"enum":        []string{"ics", "details"},
				},
			},
			"required": []string{"course", "format"},
		},
	},
}
