package main

import (
	"errors"
	"net/http"
)

func (me *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": "development",
			"version":     "1.0.0",
		},
	}

	err := me.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		me.serverErrorResponse(w, r, err)
		return
	}
}

func (me *application) searchByTagHandler(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	tag := queryValues.Get("tag")
	if tag == "" {
		me.badRequestResponse(w, r, errors.New("tag params is required"))
		return
	}

	posts, err := me.FindPostsByTag(tag)
	if err != nil {
		me.serverErrorResponse(w, r, err)
		return
	}

	err = me.writeJSON(w, http.StatusOK, envelope{"posts": posts}, nil)
	if err != nil {
		me.serverErrorResponse(w, r, err)
		return
	}

}

func (me *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := me.readIDParam(r)
	if err != nil {
		me.badRequestResponse(w, r, err)
		return
	}

	post, err := me.GetAndCachePost(int(id))

	if errors.Is(ErrNotFound, err) {
		me.notFoundResponse(w, r)
		return
	}

	err = me.writeJSON(w, http.StatusOK, envelope{
		"post": post,
	}, nil)

	if err != nil {
		me.serverErrorResponse(w, r, err)
		return
	}

}

func (me *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := me.readIDParam(r)
	if err != nil {
		me.badRequestResponse(w, r, err)
		return
	}

	var input struct {
		Title   string   `json:"title"`
		Content string   `json:"content"`
		Tags    []string `json:"tags"`
	}

	err = me.readJSON(w, r, &input)
	if err != nil {
		me.badRequestResponse(w, r, err)
		return
	}

	post := Post{
		Title:   input.Title,
		Content: input.Content,
		Tags:    input.Tags,
	}

	dbPost, err := me.UpdatePost(int(id), post)

	if err != nil {
		me.notFoundResponse(w, r)
		return
	}

	err = me.writeJSON(w, http.StatusOK, envelope{
		"post": dbPost,
	}, nil)

	if err != nil {
		me.serverErrorResponse(w, r, err)
		return
	}

}

func (me *application) createPostHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title   string   `json:"title"`
		Content string   `json:"content"`
		Tags    []string `json:"tags"`
	}

	err := me.readJSON(w, r, &input)
	if err != nil {
		me.badRequestResponse(w, r, err)
		return
	}

	post := Post{
		Title:   input.Title,
		Content: input.Content,
		Tags:    input.Tags,
	}

	dbPost, err := me.InsertPost(&post)

	if err != nil {
		me.serverErrorResponse(w, r, err)
		return
	}

	err = me.writeJSON(w, http.StatusCreated, envelope{"post": dbPost}, nil)
	if err != nil {
		me.serverErrorResponse(w, r, err)
	}
}
