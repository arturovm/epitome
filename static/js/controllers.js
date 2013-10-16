function SubscriptionsController($scope, $http) {
	$http.get('/api/subscriptions').error(function(data) {
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
				'Content-Type': 'application/x-www-form-urlencoded'
			}
		}).error(function(data) {
			if (angular.isObject(data) == true && data.error != undefined) {
				alertify.error("Error creating feed: " + data.error);
			}
		}).success(function(data) {
			$scope.url = "";
			alertify.success("Feed added");
			$http.get('/api/subscriptions').error(function(data) {
				alertify.error("Error updating feeds");
			}).success(function(data) {
				$scope.subs = data;
			});
		});
	}
	$scope.delete = function(id) {
		$http.delete('/api/subscriptions/' + id).error(function(data) {
			if (angular.isObject(data) == true && data.error != undefined) {
				alertify.error("Error deleting feed: " + data.error);
			}
		}).success(function(data) {
			alertify.success("Feed deleted");
			$http.get('/api/subscriptions').error(function(data) {
				alertify.error("Error updating feeds");
			}).success(function(data) {
				$scope.subs = data;
			});
		});
	}
}

function PreferencesController($scope, $http) {
	$scope.intervalOptions = [
		{value: "@every 5m", label: "5 minutes"},
		{value: "@every 15m", label: "15 minutes"},
		{value: "@every 30m", label: "30 minutes"},
		{value: "@every 1h", label: "hour"},
		{value: "@every 2h", label:"2 hours"},
		{value: "@every 8h", label: "8 hours"},
		{value: "@every 12h", label: "12 hours"},
		{value: "@every 1d", label: "day"},
		{value: "@every 2d", label: "2 days"}
	];
	$scope.refresh_rate = "@every 30m";

	$scope.newUserPermsOptions = [
		{value: 0, label: "Only administrators"},
		{value: 1, label: "All registered users"},
		{value: 2, label: "The entire Internet"}
	];
	$scope.new_user_permissions = 2;
	$http.get('/api/preferences').error(function(data) {
		if (angular.isObject(data) == true && data.error != undefined) {
			alertify.error("Error: ");
		}
	}).success(function(data) {
		$scope.refresh_rate = data.refresh_rate;
		$scope.new_user_permissions = data.new_user_permissions;
	});


	$scope.submit = function() {
		var prefs = {refresh_rate: $scope.refresh_rate, new_user_permissions: $scope.new_user_permissions};
		$http.put('/api/preferences', prefs).error(function(data) {
			if (angular.isObject(data) == true && data.error != undefined) {
				alertify.error("Error: ");
			}
		}).success(function(data) {
			alertify.success("Your preferences have been saved");
		});
	}
}

function LoginController($scope, $http) {
}

function SetupController($scope, $http, $filter) {
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
			});
		}
	}
}
