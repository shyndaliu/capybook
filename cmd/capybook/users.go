package main

// func (app *application) signUpHandler(w http.ResponseWriter, r *http.Request) {
// 	var input struct {
// 		Username string `json:"username"`
// 		Password string `json:"password"`
// 	}
// 	err := app.readJSON(w, r, &input)
// 	if err != nil {
// 		app.badRequestResponse(w, r, err)
// 		return
// 	}
// 	//hash, err:=bcrypt.GenerateFromPassword([]byte(input.Password), 10)
// 	if err != nil {
// 		app.badRequestResponse(w, r, err)
// 		return
// 	}

// }
