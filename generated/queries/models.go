// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package queries

import ()

type Filter struct {
	ID      int64
	Name    string
	Receipt string
}

type Request struct {
	ID          int64
	File        string
	RequesterID string
	Status      string
}
