package main

import (
	"context"
	"time"
	"tourbackend/internal/courses"
	"tourbackend/internal/courses/materials"
	"tourbackend/internal/courses/quizzes"
	"tourbackend/internal/feeds"
	"tourbackend/internal/utils"

	db "tourbackend/internal/database/gen"

	"github.com/google/uuid"
)

func createAdmin(q *db.Queries) error {
	ctx := context.Background()

	hash, err := utils.HashPassword("TdA26!")
	if err != nil {
		return err
	}

	_, err = q.CreateUser(ctx, db.CreateUserParams{
		FirstName: "lecturer",
		LastName:  "lecturer",
		Email:     "lecturer",
		Hash:      hash,
	})
	if err != nil {
		return nil
	}

	return nil
}

func seed(q *db.Queries, cs *courses.Service, fs *feeds.Service, ms *materials.Service, qs *quizzes.Service) error {

	ctx := context.Background()

	hash, err := utils.HashPassword("1234")
	if err != nil {
		panic(err)
	}

	_, err = q.CreateUser(ctx, db.CreateUserParams{
		FirstName: "Adminov",
		LastName:  "Adminsky",
		Email:     "ad.min@goabuc.cz",
		Hash:      hash,
	})
	if err != nil {
		return nil
	}

	//* Course 1

	now := time.Now().Unix()
	course1, err := cs.CreateCourse(db.CreateCourseParams{
		Uuid:        uuid.NewString(),
		Name:        "Pottery for Beginners",
		Description: "Intro into the wonderful world of pottery. No matter you experience you are welcome!",
		CreatedAt:   now,
		UpdatedAt:   now,
	}, ctx)
	if err != nil {
		return err
	}

	_, err = ms.CreateUrlMaterial(materials.CreateUrlMaterialRequest{
		CourseId:    course1.Uuid,
		MatType:     "url",
		Url:         "https://www.youtube.com/watch?v=FtES7Gd5gHE",
		Name:        "Pottery intruduction video",
		Description: "a short video to introduce students into the topic of pottery",
	}, uuid.NewString(), ctx)
	if err != nil {
		return err
	}

	_, err = ms.CreateUrlMaterial(materials.CreateUrlMaterialRequest{
		CourseId:    course1.Uuid,
		MatType:     "url",
		Url:         "https://www.youtube.com/watch?v=2taUjbCb3N8",
		Name:        "History of pottery",
		Description: "a longer video  going through the history of pottery",
	}, uuid.NewString(), ctx)
	if err != nil {
		return err
	}

	one := 1
	quiz1 := quizzes.Quiz{
		Uuid:          uuid.NewString(),
		Title:         "Pottery Videos Progress Assesment",
		AttemptsCount: 0,
		Questions: []quizzes.Question{
			{
				Uuid:         uuid.NewString(),
				QueType:      "singleChoice",
				Question:     "How many videos were there",
				Options:      []string{"one", "two"},
				CorrectIndex: &one,
			},
			{
				Uuid:           uuid.NewString(),
				QueType:        "multipleChoice",
				Question:       "Because of what must we suffer",
				Options:        []string{"god", "devil", "us"},
				CorrectIndices: []int{0, 1, 2},
			},
		},
	}
	_, err = qs.CreateQuiz(quiz1, course1.Uuid, ctx)
	if err != nil {
		return err
	}

	_, err = fs.CreateAutomaticPost("Pottery Course Has Been Published!", course1.Uuid, ctx)
	if err != nil {
		return err
	}

	//* Course 2

	course2, err := cs.CreateCourse(db.CreateCourseParams{
		Uuid:        uuid.NewString(),
		Name:        "Potions 101",
		Description: "Intro into potion making, fast-paced course for serious sorcerers only",
		CreatedAt:   now,
		UpdatedAt:   now,
	}, ctx)
	if err != nil {
		return err
	}

	_, err = ms.CreateUrlMaterial(materials.CreateUrlMaterialRequest{
		CourseId:    course2.Uuid,
		MatType:     "url",
		Url:         "https://www.youtube.com/watch?v=ND-h-Qxym1M",
		Name:        "Potions basics",
		Description: "a short video to introduce students into the topic of pottery",
	}, uuid.NewString(), ctx)

	_, err = ms.CreateUrlMaterial(materials.CreateUrlMaterialRequest{
		CourseId:    course2.Uuid,
		MatType:     "url",
		Url:         "https://www.youtube.com/watch?v=K0sIBsp-d6A",
		Name:        "Potions in popular culture",
		Description: "blah blah lorem ipsum etcetera casius belli",
	}, uuid.NewString(), ctx)

	//* Course 3

	course3, err := cs.CreateCourse(db.CreateCourseParams{
		Uuid:        uuid.NewString(),
		Name:        "Zebra Riding Advanced",
		Description: "A guide to advanced zebra riding techniques, must already own a zebra",
		CreatedAt:   now,
		UpdatedAt:   now,
	}, ctx)
	if err != nil {
		return err
	}

	_, err = ms.CreateUrlMaterial(materials.CreateUrlMaterialRequest{
		CourseId:    course3.Uuid,
		MatType:     "url",
		Url:         "https://www.youtube.com/watch?v=Ph8Vag9VxRU",
		Description: "What will you learn",
	}, uuid.NewString(), ctx)

	_, err = ms.CreateUrlMaterial(materials.CreateUrlMaterialRequest{
		CourseId:    course3.Uuid,
		MatType:     "url",
		Url:         "https://www.youtube.com/watch?v=grjZPfCH6bs",
		Description: "What i wish i had for dinner today",
	}, uuid.NewString(), ctx)

	return nil
}
