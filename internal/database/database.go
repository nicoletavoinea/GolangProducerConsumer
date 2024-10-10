package database

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	//"sync"
	proto "github.com/nicoletavoinea/GolangProducerConsumer/api/proto/task"
	sqlc "github.com/nicoletavoinea/GolangProducerConsumer/internal/database/sqlc"
	zl "github.com/rs/zerolog/log"
)

var Queries *sqlc.Queries
var Db *sql.DB

//var mutexAdd sync.Mutex

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "postgres"
)

func OpenDatabase() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Failed to open postgres database: %v", err)
	}

	if err := CreateType(db, "../../internal/database/queries/createType.sql"); err != nil {
		log.Fatalf("Error in CreateType file: %v", err)
	}

	if err := RunSchema(db, "../../internal/database/queries/schema.sql"); err != nil {
		log.Fatalf("Failed to run schema: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Ping failed: %v", err)
	}

	fmt.Println("Successfully connected to PostgreSQL database!")

	Queries = sqlc.New(db)
	return db
}

func CreateType(db *sql.DB, createTypeFilePath string) error {
	createtype, err := ioutil.ReadFile(createTypeFilePath)
	if err != nil {
		return err
	}
	_, err = db.Exec(string(createtype))

	if err != nil && !strings.Contains(err.Error(), "already exists") { // Check if type exists & if already exists do nothing
		return err // Return error if it fails for another reason
	}
	return nil
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
		log.Fatal("Error closing database connection:", err)
	}
}

func AddTaskToDatabase(task *proto.Task, db *sqlc.Queries) error {
	//mutexAdd.Lock()
	taskData, err := db.AddTask(context.Background(), sqlc.AddTaskParams{
		Type:  task.TaskType,
		Value: task.TaskValue,
	})

	//mutexAdd.Unlock()
	if err != nil {
		log.Printf("Error inserting task: %v\n", err)
		return err
	}

	task.TaskId = taskData.ID
	task.TaskCreationTime = taskData.Creationtime
	task.TaskLastUpdateTime = task.TaskCreationTime
	//log.Printf("Added: %v\n", task)
	zl.Info().Int32("TaskId", task.TaskId).Int32("TaskType", task.TaskType).Int32("TaskValue", task.TaskValue).Str("State", string(sqlc.TaskStateRECEIVED)).Msg("Added to database by producer")
	return nil
}

func UpdateTaskState(taskID int32, status proto.TaskState) (sqlc.Task, error) {
	var state sqlc.TaskState
	if status == proto.TaskState_PROCESSING {
		state = sqlc.TaskStatePROCESSING
	} else if status == proto.TaskState_DONE {
		state = sqlc.TaskStateDONE
	}

	updatedTask, err := Queries.UpdateTask(context.Background(), sqlc.UpdateTaskParams{
		ID:      taskID,
		Column2: state,
	})
	if err != nil {
		log.Printf("Error updating task: %v\n", err)
		return updatedTask, err
	}
	//log.Printf("Updated: %v\n", updatedTask)
	zl.Info().Int32("TaskId", updatedTask.ID).Int32("TaskType", updatedTask.Type).Int32("TaskValue", updatedTask.Value).Str("State", string(updatedTask.State)).Msg("Database updated by consumer")
	return updatedTask, nil
}

func GetNumberOfDoneTasks() float64 { //retrieve from database the number of tasks that are in the DONE state

	doneTasks, err := Queries.GetNumberOfTasks(context.Background(), sqlc.TaskStateDONE)
	if err != nil {
		log.Printf("Error getting the number of done tasks task: %v\n", err)
		return 0
	}
	return float64(doneTasks)
}

func GetNumberOfProcessingTasks() float64 { //retrieve from database the number of tasks that are in the PROCESSING state
	processingTasks, err := Queries.GetNumberOfTasks(context.Background(), sqlc.TaskStatePROCESSING)
	if err != nil {
		log.Printf("Error getting the number of processing tasks task: %v\n", err)
		return 0
	}
	return float64(processingTasks)
}

func GetNumberOfTasksByType() [10]int { //retrieve from database the number of tasks of each type
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

func GetValueOfTasksByType() [10]int64 { //retrieve from database the sum of the values of each task type
	values := [10]int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	rows, err := Queries.GetValueOfTasksByType(context.Background())
	if err != nil {
		log.Printf("Error getting the value of tasks by type: %v\n", err)
		return values
	}
	for i := 0; i < len(rows); i++ {
		values[rows[i].Type] = rows[i].ValuesSum
	}

	return values
}
