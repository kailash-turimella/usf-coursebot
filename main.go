package main

const chromadbHost = "http://localhost:8000"
const fileName = "Fall 2024 Class Schedule 08082024.csv"
const useTerminal = false
const printResponse = false
const numRelevantCoursesToRetrieve = 30
const systemPrompt = `Use the provided list of USF Fall 2024 courses to answer questions.

Here are the columns in the CSV file (containing all course data) and an example row:
SUBJ,CRSE NUM,SEC,CRN,Schedule Type Code,Campus Code,Title Short Desc,Instruction Mode Desc,Meeting Type Codes,Meet Days,Begin Time,End Time,Meet Start,Meet End,BLDG,RM,Actual Enrollment,Primary Instructor First Name,Primary Instructor Last Name,Primary Instructor Email,College
ADVT,401,01,40328,FWK,M,Advertising Internship,In-Person,IP,S,1145,1525,8/20/24,11/30/24,LM,141A,5,David,McGrane,dmcgrane@usfca.edu,LA

For questions that makes you list things, please list them in a neat format.

Only call tools if the user's question requires exact data returned by a specific tool. It's okay to not call a tool. You may call multiple tools if necessary but NO UNECESSARY TOOL CALLS`

const getJSONPromptTemplate = `Given a natural language question about USF Fall 2024 courses, identify the most relevant CSV column to search, and the value to search for.

Return a JSON object FOR EXAMPLE:
{
  "column": "Primary Instructor Last Name",
  "value": "Peterson"
}

Only return the JSON. No explanation or extra text.

QUESTION : `

func main() {
	collectionData := setUpCollections(fileName)
	if collectionData.collection != nil {
		startChatbot(collectionData)
	}
}

// RUN THE CODE USING : go run $(find . -name "*.go" ! -name "*_test.go")
