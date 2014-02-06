function SubscriptionsController($scope, $document, $http, $cookieStore, $location) {
	$http.get('/api/subscriptions', {
		headers: {
			'X-Session-Token': $cookieStore.get('user').session_token
		}
	}).error(function(data) {
		alertify.error("Error updating feeds");
	}).success(function(data) {
		$scope.subs = data;
	});
	$scope.submit = function() {
		alertify.log("Adding feed...", "", 1500);
		var url = $scope.url;
		if (url == undefined || url == "") {
			alertify.error("Error: Please fill in the URL field");
			return;
		}
		$http.post('/api/subscriptions', "url=" + url, {
			headers: {
				'Content-Type': 'application/x-www-form-urlencoded',
				'X-Session-Token': $cookieStore.get('user').session_token
			}
		}).error(function(data) {
			if (angular.isObject(data) == true && data.error != undefined) {
				alertify.error("Error creating feed: " + data.error);
			}
		}).success(function(data) {
			$scope.url = "";
			alertify.success("Feed added");
			$http.get('/api/subscriptions', {
				headers: {
					'X-Session-Token': $cookieStore.get('user').session_token
				}
			}).error(function(data) {
				alertify.error("Error updating feeds");
			}).success(function(data) {
				$scope.subs = data;
			});
		});
	};
	/*$scope.submitFile = function() {
		var payload = new FormData();
		payload.append("subscriptions", document.getElementById("subs-file").files[0]);
		$http.post('/api/subscriptions', payload, {
			headers: {
				'Content-Type': undefined,
				'X-Session-Token': $cookieStore.get('user').session_token
			},
			transformRequest: function(data) {
				return data
			}
		}).error(function(data) {
			if (angular.isObject(data) == true && data.error != undefined) {
				alertify.error("Error uploading file: " + data.error);
			}
		}).success(function(data) {
			alertify.success("Feeds added");
			$http.get('/api/subscriptions', {
				headers: {
					'x-session-token': $cookieStore.get('user').session_token
				}
			}).error(function(data) {
				alertify.error("Error updating feeds");
			}).success(function(data) {
				$scope.subs = data;
			});
		});
	};*/
	$scope.delete = function(id) {
		$http.delete('/api/subscriptions/' + id, {
			headers: {
				'X-Session-Token': $cookieStore.get('user').session_token
			}
		}).error(function(data) {
			if (angular.isObject(data) == true && data.error != undefined) {
				alertify.error("Error deleting feed: " + data.error);
			}
		}).success(function(data) {
			alertify.success("Feed deleted");
			$http.get('/api/subscriptions', {
				headers: {
					'x-session-token': $cookieStore.get('user').session_token
				}
			}).error(function(data) {
				alertify.error("Error updating feeds");
			}).success(function(data) {
				$scope.subs = data;
			});
		});
	};
	$scope.logout = function() {
		$http.delete('/api/auth/sessions/' + $cookieStore.get('user').session_token, {
			headers: {
				'X-Session-Token': $cookieStore.get('user').session_token
			}
		}).error(function(data) {
			if (angular.isObject(data) == true && data.error != undefined) {
				alertify.error("Error logging you out: " + data.error);
			}
		}).success(function(data) {
			$cookieStore.remove('user');
			$location.path('/login');
			alertify.success("You are now logged out");
		});
	};
}

function PreferencesController($scope, $http, $cookieStore) {
	$scope.intervalOptions = [
		{value: "15m", label: "15 minutes"},
		{value: "30m", label: "30 minutes"},
		{value: "1h", label: "hour"},
		{value: "2h", label:"2 hours"},
		{value: "8h", label: "8 hours"},
		{value: "12h", label: "12 hours"},
		{value: "1d", label: "day"},
		{value: "2d", label: "2 days"}
	];
	$scope.refresh_rate = "@every 30m";

	$scope.newUserPermsOptions = [
		{value: 0, label: "Only administrators"},
		{value: 1, label: "All registered users"},
		{value: 2, label: "The entire Internet"}
	];
	$scope.new_user_permissions = 2;
	$http.get('/api/preferences', {
		headers: {
			'X-Session-Token': $cookieStore.get('user').session_token
		}
	}).error(function(data) {
		if (angular.isObject(data) == true && data.error != undefined) {
			alertify.error("Error: " + data.error);
		}
	}).success(function(data) {
		$scope.refresh_rate = data.refresh_rate;
		$scope.new_user_permissions = data.new_user_permissions;
	});


	$scope.submit = function() {
		var prefs = {refresh_rate: $scope.refresh_rate, new_user_permissions: $scope.new_user_permissions};
		$http.put('/api/preferences', prefs, {
			headers: {
				'X-Session-Token': $cookieStore.get('user').session_token
			}
		}).error(function(data) {
			if (angular.isObject(data) == true && data.error != undefined) {
				alertify.error("Error: " + data.error);
			}
		}).success(function(data) {
			alertify.success("Your preferences have been saved");
		});
	};
	$scope.logout = function() {
		$http.delete('/api/auth/sessions/' + $cookieStore.get('user').session_token, {
			headers: {
				'X-Session-Token': $cookieStore.get('user').session_token
			}
		}).error(function(data) {
			if (angular.isObject(data) == true && data.error != undefined) {
				alertify.error("Error logging you out: " + data.error);
			}
		}).success(function(data) {
			$cookieStore.remove('user');
			$location.path('/login');
			alertify.success("You are now logged out");
		});
	};
}

function LoginController($scope, $http, $filter, $cookieStore, $location) {
	$scope.submit = function() {
		if ($scope.username == undefined || $scope.password == undefined || $scope.username == "" || $scope.password == "") {
			alertify.error("Please fill in all the fields.");
		} else {
			var user = $scope.username;
			var pass = SparkMD5.hash($filter("lowercase")($scope.username) + ":" + $scope.password);
			$http.post('/api/auth/sessions', "username=" + user + "&password=" + pass + "&appname=Pond Admin", {
				headers: {
					'Content-Type': 'application/x-www-form-urlencoded'
				}
			}).error(function(data) {
				if (angular.isObject(data) == true && data.error != undefined) {
					alertify.error("Error: " + data.error);
				}
			}).success(function(data) {
				if (angular.isObject(data) == true && data.session_token != undefined) {
					$cookieStore.put('user', {
						username: user,
						session_token: data.session_token
					});
					$location.path('/subscriptions');
				}
			});
		}
	};
}

function SetupController($scope, $http, $filter) {
	$scope.accountCreated = false;
	$scope.submit = function () {
		if ($scope.username == undefined || $scope.password == undefined || $scope.username == "" || $scope.password == "") {
			alertify.error("Please fill in all the fields.");
		} else {
			var user = $scope.username;
			var pass = SparkMD5.hash($filter("lowercase")($scope.username) + ":" + $scope.password);
			$http.post('/api/users', "username=" + user + "&password=" + pass + "&role=admin", {
				headers: {
					'Content-Type': 'application/x-www-form-urlencoded'
				}
			}).error(function(data) {
				if (angular.isObject(data) == true && data.error != undefined) {
					alertify.error("Error: " + data.error);
				}
			}).success(function(data) {
				$scope.accountCreated = true;
				alertify.success("Your account has been created successully");
			});
		}
	};
}
