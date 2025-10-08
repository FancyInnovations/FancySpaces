package handler

import (
	"net/http"

	"github.com/OliverSchlueter/goutils/problems"
)

func (h *Handler) handleMembers(w http.ResponseWriter, r *http.Request) {
	// TODO implement

	// POST: add member (only owner can add)
	// GET: list members (no auth required)

	problems.NotImplemented().WriteToHTTP(w)
}

func (h *Handler) handleMember(w http.ResponseWriter, r *http.Request) {
	// TODO implement

	// DELETE: remove member (only owner can remove and the user themselves)
	// PUT: change member role (only owner change roles)

	problems.NotImplemented().WriteToHTTP(w)
}
