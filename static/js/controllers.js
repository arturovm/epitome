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
}
