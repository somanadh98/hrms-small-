package routes

import (
    "github.com/gorilla/mux"
    "gorm.io/gorm"
)

func Register(r *mux.Router, db *gorm.DB) {
    registerAuthRoutes(r, db)
    registerEmployeeRoutes(r, db)
    registerAttendanceRoutes(r, db)
    registerLeaveRoutes(r, db)
}


