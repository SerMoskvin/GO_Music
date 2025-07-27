package app



r := mux.NewRouter()
userHandler := api.NewUserHandler(userManager, sessionStore)

r.HandleFunc("/users", userHandler.CreateUserHandler).Methods("POST")
r.HandleFunc("/users", userHandler.GetUsersByIDsHandler).Methods("GET") // с параметром ids
r.HandleFunc("/users/{id:[0-9]+}", userHandler.GetUserHandler).Methods("GET")
r.HandleFunc("/users/{id:[0-9]+}", userHandler.UpdateUserHandler).Methods("PUT")
r.HandleFunc("/users/{id:[0-9]+}", userHandler.DeleteUserHandler).Methods("DELETE")
