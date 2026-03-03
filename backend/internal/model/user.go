package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username     string             `bson:"username" json:"username"`
	Password     string             `bson:"password,omitempty" json:"-"`
	Email        string             `bson:"email" json:"email"`
	Gender       string             `bson:"gender" json:"gender"`
	Age          int                `bson:"age" json:"age"`
	ProjectState ProjectState       `bson:"project_state,omitempty" json:"project_state,omitempty"`
	CreatedAt    int64              `bson:"created_at" json:"created_at"`
	LastLoginAt  int64              `bson:"last_login_at,omitempty" json:"last_login_at,omitempty"`
}

type ProjectItem struct {
	ID         int    `bson:"id" json:"id"`
	Name       string `bson:"name" json:"name"`
	Owner      string `bson:"owner" json:"owner"`
	Center     string `bson:"center" json:"center"`
	State      string `bson:"state" json:"state"`
	StateClass string `bson:"stateClass" json:"stateClass"`
	Desc       string `bson:"desc" json:"desc"`
	Date       string `bson:"date" json:"date"`
	Members    int    `bson:"members" json:"members"`
}

type TodoItem struct {
	Key         string `bson:"key" json:"key"`
	ProjectID   int    `bson:"projectId" json:"projectId"`
	ProjectName string `bson:"projectName" json:"projectName"`
	TaskName    string `bson:"taskName" json:"taskName"`
}

type ProjectState struct {
	RecentProjects []ProjectItem `bson:"recentProjects" json:"recentProjects"`
	TodoItems      []TodoItem    `bson:"todoItems" json:"todoItems"`
	NextProjectID  int           `bson:"nextProjectId" json:"nextProjectId"`
}
