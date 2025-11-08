package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"
)

type NoteController struct {
	baseController BaseController[schema.Note]
}

func NewNoteController() NoteController {
	return NoteController{
		baseController: BaseController[schema.Note]{
			collectionName: "notes",
			displayName:    "Note",
		},
	}
}

func validateNoteCreator(c *gin.Context) {

}

func (nc NoteController) Create(c *gin.Context) {
	nc.baseController.Create(c)
}

func (nc NoteController) Delete(c *gin.Context) {
	nc.baseController.Delete(c)
}

func (nc NoteController) Update(c *gin.Context) {
	nc.baseController.Update(c)
}

func (nc NoteController) RetrieveAll(c *gin.Context) {
	nc.baseController.RetrieveAll(c)
}

func (nc NoteController) RetrieveOne(c *gin.Context) {
	nc.baseController.RetrieveOne(c)
}
