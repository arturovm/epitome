angular.module('pond', ['ngRoute', 'ngCookies']).config(['$routeProvider', function($routeProvider) {
	$routeProvider.when('/subscriptions', {
		templateUrl: 'partials/subscriptions_view.html',
		controller: 'SubscriptionsController'
	}).when('/preferences', {
		templateUrl: 'partials/preferences_view.html',
		controller: 'PreferencesController'
	}).when('/login', {
		templateUrl: 'partials/login_view.html',
		controller: 'LoginController'
	}).otherwise({
		redirectTo: '/subscriptions'
	});
}]).run(function($rootScope, $location, $cookieStore) {
	$rootScope.$on("$routeChangeStart", function(event, next, last) {
		if (!angular.isDefined($cookieStore.get("user"))) {
			if (next.templateUrl != '/partials/login_view.html') {
				$location.path("/login");
			}
	}
	});
});
