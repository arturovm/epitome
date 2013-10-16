package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	//"log"
)

func GetSetup(w http.ResponseWriter, req *http.Request) {
	DB, _ := sql.Open("sqlite3", ExePath+"/db.db")
	rows, err := DB.Query("select * from users")
	DB.Close()
	if err != nil {
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`Error: Could not connect to database`))
		return
	}
	users := make([]User, 0)
	for rows.Next() {
		var user User
		rows.Scan(&user.Username, &user.PasswordHash)
		users = append(users, user)
	}
	rows.Close()
	if len(users) > 0 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`You are not authorized to view this page`))
		return
	}
	var HTML = `
	<!DOCTYPE html>
	<html lang="en" ng-app="pond">
		<head>
			<title>Pond Setup</title>
			<meta charset="utf-8">
			<link rel="stylesheet" href="../css/base.css">
			<link rel="stylesheet" href="../css/skeleton.css">
			<link rel="stylesheet" href="../css/layout.css">
			<link rel="stylesheet" href="http://yui.yahooapis.com/pure/0.3.0/pure-min.css">		
			<link rel="stylesheet" href="../bower_components/alertify/themes/alertify.core.css">
			<link rel="stylesheet" href="../bower_components/alertify/themes/alertify.default.css">
			<link rel="stylesheet" href="../css/main.css">
		</head>
		<body ng-controller="SetupController">
			<div class="container">
				<div class="four columns offset-by-six">
					<form ng-submit="submit()" class="pure-form pure-form-stacked">
						<!--<legend>Create a new account</legend>
						<label for="username">Username:</label>
						<input type="text" id="username" ng-model="username">
						<label for="password">Password:</label>
						<input type="password" id="password" ng-model="password">
						<a type="submit" class="pure-button pure-button-primary" ng-click="submit()">Create</a>-->

						<label for="username">Username</label>
						<input type="text" id="username" ng-model="username">
						<label for="password">Password</label>
						<input type="password" id="password" ng-model="password">
						<button type="submit" class="pure-button pure-button-primary" value="Create">Create</button>
					</form>
				</div>
			</div>
		<script async src="http://use.edgefonts.net/source-sans-pro;droid-serif;merriweather;droid-sans.js"></script>
		<script src="../bower_components/angular/angular.min.js"></script>
		<script src="../bower_components/angular-route/angular-route.min.js"></script>
		<script src="../bower_components/angular-cookies/angular-cookies.min.js"></script>
		<script src="../bower_components/alertify/alertify.min.js"></script>
		<script src="../bower_components/spark-md5/spark-md5.min.js"></script>
		<script src="../js/controllers.js"></script>
		<script src="../js/app.js"></script>
		</body>
	</html>
	`
	w.Header().Set("content-type", "text/html")
	w.Write([]byte(HTML))
}
