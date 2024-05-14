package usecase

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	dbModel "simple-users-tasks-service/db/model"
	handlerModel "simple-users-tasks-service/handler/model"

	"gorm.io/gorm"
)

type Usecase struct {
	db *gorm.DB
}

func NewUsecase(db *gorm.DB) Usecase {
	return Usecase{
		db: db,
	}
}

// User usecases
func (u *Usecase) ProcessRegisterUser(reqBody handlerModel.RegisterUserRequest) (statusCode int, err error) {
	hasher := sha256.New()
	hasher.Write([]byte(reqBody.Password))
	hashedPassword := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	if err = u.db.Create(&dbModel.User{
		Email:    reqBody.Email,
		Password: string(hashedPassword),
		Name:     reqBody.Name,
	}).Error; err != nil {
		log.Println(err)
		return http.StatusInternalServerError, fmt.Errorf("error inserting new user to DB")
	}

	return http.StatusCreated, nil
}

func (u *Usecase) ProcessGetAllUsers() (data []dbModel.User, statusCode int, err error) {
	rows, err := u.db.Select("id,name").Model(&dbModel.User{}).Rows()
	if err != nil {
		defer rows.Close()
		log.Println(err)
		return nil, http.StatusInternalServerError, fmt.Errorf("unable to get user data from DB")
	}
	defer rows.Close()

	var users []dbModel.User

	for rows.Next() {
		err = u.db.ScanRows(rows, &users)
		if err != nil {
			log.Println(err)
			return nil, http.StatusInternalServerError, fmt.Errorf("unable to parse user data")
		}

	}

	return users, http.StatusOK, nil
}

func (u *Usecase) ProcessGetUserDetails(id int) (data dbModel.User, statusCode int, err error) {
	var user = dbModel.User{ID: id}
	result := u.db.Select("id,name").First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		return dbModel.User{}, http.StatusNotFound, fmt.Errorf("user not found")
	}
	if result.Error != nil {
		log.Println(err)
		return dbModel.User{}, http.StatusInternalServerError, fmt.Errorf("unable to parse user data")

	}
	return user, http.StatusOK, nil
}

func (u *Usecase) ProcessUpdateUserData(id int, reqBody handlerModel.UpdateUserRequest) (statusCode int, err error) {
	var user = dbModel.User{ID: id}
	result := u.db.First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		return http.StatusNotFound, fmt.Errorf("user not found")
	}
	if result.Error != nil {
		log.Println(err)
		return http.StatusInternalServerError, fmt.Errorf("unable to parse user data")
	}

	updateFlag := false
	if reqBody.Email != user.Email {
		updateFlag = true
		user.Email = reqBody.Email
	}

	if reqBody.Name != user.Name {
		updateFlag = true
		user.Name = reqBody.Name
	}

	if updateFlag {
		updateProcess := result.Updates(user)
		if updateProcess.Error != nil {
			return http.StatusInternalServerError, fmt.Errorf("unable to update user data")
		}
	} else {
		return http.StatusBadRequest, fmt.Errorf("no changes made")

	}

	return http.StatusNoContent, nil
}

func (u *Usecase) ProcessDeleteUserData(id int) (statusCode int, err error) {
	var user = dbModel.User{ID: id}
	result := u.db.First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		return http.StatusNotFound, fmt.Errorf("user not found")
	}

	if result.Error != nil {
		log.Println(err)
		return http.StatusInternalServerError, fmt.Errorf("unable to parse user data")
	}

	deleteProcess := u.db.Delete(&dbModel.User{}, id)
	if deleteProcess.Error != nil {
		log.Println(err)
		return http.StatusInternalServerError, fmt.Errorf("error while deleting user data")
	}
	return http.StatusNoContent, nil
}

// Task usecases
func (u *Usecase) ProcessCreateTask(reqBody handlerModel.CreateTaskRequest) (statusCode int, err error) {
	if err = u.db.Create(&dbModel.Task{
		UserID:      reqBody.UserID,
		Title:       reqBody.Title,
		Description: reqBody.Description,
	}).Error; err != nil {
		log.Println(err)
		return http.StatusInternalServerError, fmt.Errorf("error inserting new task to DB")
	}

	return http.StatusCreated, nil
}

func (u *Usecase) ProcessGetAllTasks() (data []dbModel.Task, statusCode int, err error) {
	rows, err := u.db.Select("id,user_id,title,description").Model(&dbModel.Task{}).Rows()
	if err != nil {
		defer rows.Close()
		log.Println(err)
		return nil, http.StatusInternalServerError, fmt.Errorf("unable to get task data from DB")
	}
	defer rows.Close()

	var tasks []dbModel.Task

	for rows.Next() {
		err = u.db.ScanRows(rows, &tasks)
		if err != nil {
			log.Println(err)
			return nil, http.StatusInternalServerError, fmt.Errorf("unable to parse task data")
		}

	}

	return tasks, http.StatusOK, nil
}

func (u *Usecase) ProcessGetTaskDetails(id int) (data dbModel.Task, statusCode int, err error) {
	var task = dbModel.Task{ID: id}
	result := u.db.First(&task)
	if result.Error == gorm.ErrRecordNotFound {
		return dbModel.Task{}, http.StatusNotFound, fmt.Errorf("task not found")
	}
	if result.Error != nil {
		log.Println(err)
		return dbModel.Task{}, http.StatusInternalServerError, fmt.Errorf("unable to parse task data")

	}
	return task, http.StatusOK, nil
}

func (u *Usecase) ProcessUpdateTask(id int, reqBody handlerModel.UpdateTaskRequest) (statusCode int, err error) {
	var task = dbModel.Task{ID: id}
	result := u.db.First(&task)
	if result.Error == gorm.ErrRecordNotFound {
		return http.StatusNotFound, fmt.Errorf("task not found")
	}
	if result.Error != nil {
		log.Println(err)
		return http.StatusInternalServerError, fmt.Errorf("unable to parse task data")
	}

	updateFlag := false
	if reqBody.Title != task.Title {
		updateFlag = true
		task.Title = reqBody.Title
	}

	if reqBody.Description != task.Description {
		updateFlag = true
		task.Description = reqBody.Description
	}

	if reqBody.Status != task.Status {
		updateFlag = true
		task.Status = reqBody.Status
	}

	if reqBody.UserID != task.UserID {
		updateFlag = true
		task.UserID = reqBody.UserID
	}

	if updateFlag {
		updateProcess := result.Updates(task)
		if updateProcess.Error != nil {
			return http.StatusInternalServerError, fmt.Errorf("unable to update task data")
		}
	} else {
		return http.StatusBadRequest, fmt.Errorf("no changes made")

	}

	return http.StatusNoContent, nil
}

func (u *Usecase) ProcessDeleteTask(id int) (statusCode int, err error) {
	var task = dbModel.Task{ID: id}
	result := u.db.First(&task)
	if result.Error == gorm.ErrRecordNotFound {
		return http.StatusNotFound, fmt.Errorf("task not found")
	}

	if result.Error != nil {
		log.Println(err)
		return http.StatusInternalServerError, fmt.Errorf("unable to parse task data")
	}

	deleteProcess := u.db.Delete(&dbModel.Task{}, id)
	if deleteProcess.Error != nil {
		log.Println(err)
		return http.StatusInternalServerError, fmt.Errorf("error while deleting task data")
	}
	return http.StatusNoContent, nil
}
