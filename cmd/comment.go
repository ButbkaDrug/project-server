package main

import (
	"lp/pkg/comments"
	"lp/pkg/lyrics"
	"net/http"
	"os"
	"strconv"
)

// this endpoint will save a new comment
func (app *app) commentHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "failed to parse form values", http.StatusBadRequest)
		app.logger.Error("cmd/comment: failed to parse form values", "error", err.Error())
		return
	}

	iid := r.Form.Get("itemID")
	sid := r.Form.Get("sectionID")
	text := r.Form.Get("text")
	commentType := r.Form.Get("type")

	if iid == "" || sid == "" {
		http.Error(w, "itemID and sectionID should be set", http.StatusBadRequest)
		app.logger.Error("cmd/comment: itemID or sectionID not set", "itemID", iid, "sectionID", sid)
		return
	}

	itemID, iErr := strconv.Atoi(iid)
	sectionID, sErr := strconv.Atoi(sid)

	if iErr != nil || sErr != nil {
		http.Error(w, "itemID and sectionID should be valid numbers", http.StatusInternalServerError)
		app.logger.Error("cmd/comment: itemID or sectionID is invalid", "iError", iErr.Error(), "sError", sErr.Error())
		return
	}

	comment := struct {
		ItemID    int
		SectionID int
		Text      string
	}{
		itemID,
		sectionID,
		text,
	}

	switch r.Method {
	case "GET":
		err := app.renderer.RenderComponent(w, comment, "add-comment")
		if err != nil {
			http.Error(w, "failed to render comment form", http.StatusInternalServerError)
			app.logger.Error("cmd/comment: render(add-comment) failed", "error", err)
			return
		}
	case "PUT":

		if itemID >= len(app.reaper.Comments) {
			http.Error(w, "itemID is too big", http.StatusInternalServerError)
			app.logger.Error("cmd/comment: itemID is too big", "itemID", itemID, "len(Comments)", len(app.reaper.Comments))
			return

		}

		if text == "" {
			http.Error(w, "comment cannot be emtpy", http.StatusBadRequest)
			app.logger.Error("cmd/comment: comment text is empty!")
			return
		}

		item := app.reaper.Comments[itemID]

		id := -1

		for i := range len(item.Sections) {

			if item.Sections[i].Number == sectionID {
				id = i
				break
			}

		}

		if id < 0 {
			item.Sections = append(item.Sections, lyrics.Section{})
			id = len(item.Sections) - 1
		}

		if item.Sections[id].Empty() {
			item.Sections[id].Name = "Notes"
			item.Sections[id].Number = sectionID

		}

		newComment := comments.NewComment()

		newComment.Text = comment.Text
		newComment.Type = comments.GetCommentType(commentType)

		item.Sections[id].Comments = append(item.Sections[id].Comments, newComment)

		//FIXME: there is potential race condition or a lock situation if two rootens will try to write to a file. At this stage is very, very unlickly. But needs to be fixed
		file, err := os.Create(item.Path)

		if err != nil {
			http.Error(w, "faild to open a file.", http.StatusInternalServerError)
			app.logger.Error("cmd/comment: failed to open/create a file", "path", item.Path, "error", err.Error())
			return
		}
		defer file.Close()
		if _, err := item.Write(file); err != nil {
			http.Error(w, "failed to save changes to a file", http.StatusInternalServerError)
			app.logger.Error("cmd/comment: failed to write changes to a file", "error", err.Error())
			return
		}

		if err := app.renderer.RenderComponent(w, comment.Text, "primer-comment"); err != nil {
			http.Error(w, "failed to render new comment", http.StatusInternalServerError)
			app.logger.Error("cmd/comment: render(primer-comment) failed", "error", err)
			return
		}
		if err := app.renderer.RenderComponent(w, comment, "add-comment-button"); err != nil {
			http.Error(w, "failed to render new comment", http.StatusInternalServerError)
			app.logger.Error("cmd/comment: render(add-comment-button) failed", "error", err)
			return
		}

	}

}
