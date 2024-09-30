package functions

import (
	"context"
	"database/sql"
	"io/ioutil"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	database "github.com/nicoletavoinea/GolangProducerConsumer/database/sqlc"
)

var Queries *database.Queries
var Db *sql.DB
var mutexAdd sync.Mutex

func OpenDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "../database/tasks.db")
	if err != nil {
		log.Fatal(err)
	}

	if err := RunSchema(db, "../database/schema.sql"); err != nil {
		log.Fatalf("Failed to run schema: %v", err)
	}

	Queries = database.New(db)
	return db
}

func RunSchema(db *sql.DB, schemaFilePath string) error {
	schema, err := ioutil.ReadFile(schemaFilePath)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(schema))
	return err
}

func CloseDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Fatalf("Error closing database: %v", err)
	}
}

func AddTaskToDatabase(task task, db *database.Queries) (task, error) {
	mutexAdd.Lock()
	taskData, err := db.AddTask(context.Background(), database.AddTaskParams{
		Param1: int64(task.TaskType),
		Param2: int64(task.TaskValue),
	})

	mutexAdd.Unlock()
	if err != nil {
		log.Printf("Error inserting task: %v\n", err)
		return task, err
	}
	task.TaskId = int32(taskData.ID)
	task.TaskCreationTime = taskData.Creationtime
	task.TaskLastUpdateTime = task.TaskCreationTime
	return task, nil
}

func UpdateTaskState(taskID int32, status StatusCode) (database.Task, error) {
	state := ""
	if status == PROCESSING {
		state = "PROCESSING"
	} else if status == DONE {
		state = "DONE"
	}

	updatedTask, err := Queries.UpdateTask(context.Background(), database.UpdateTaskParams{
		Param1: int64(taskID),
		Param2: state,
	})
	if err != nil {
		log.Printf("Error updating task: %v\n", err)
		return updatedTask, err
	}
	return updatedTask, nil
}

func getNumberOfDoneTasks() float64 {

	doneTasks, err := Queries.GetNumberOfTasks(context.Background(), "DONE")
	if err != nil {
		log.Printf("Error getting the number of done tasks task: %v\n", err)
		return 0
	}
	return float64(doneTasks)
}

func getNumberOfProcessingTasks() float64 {
	processingTasks, err := Queries.GetNumberOfTasks(context.Background(), "PROCESSING")
	if err != nil {
		log.Printf("Error getting the number of done tasks task: %v\n", err)
		return 0
	}
	return float64(processingTasks)
}

func getNumberOfTasksByType() [10]int {
	values := [10]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	rows, err := Queries.GetNumberOfTasksByType(context.Background())
	if err != nil {
		log.Printf("Error getting the number of tasks by type: %v\n", err)
		return values
	}

	for i := 0; i < len(rows); i++ {
		values[rows[i].Type] = int(rows[i].TaskCount)
	}

	return values
}

func getValueOfTasksByType() [10]float64 {
	values := [10]float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	rows, err := Queries.GetValueOfTasksByType(context.Background())
	if err != nil {
		log.Printf("Error getting the value of tasks by type: %v\n", err)
		return values
	}

	for i := 0; i < len(rows); i++ {
		values[rows[i].Type] = rows[i].ValuesSum.Float64
	}

	return values
}
