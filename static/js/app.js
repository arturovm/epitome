angular.module('pond', ['ngRoute']).config(['$routeProvider', function($routeProvider) {
	$routeProvider.when('/subscriptions', {
		templateUrl: 'partials/subscriptions_view.html',
		controller: 'SubscriptionsController'
	}).when('/preferences', {
		templateUrl: 'partials/preferences_view.html',
		controller: 'PreferencesController'
	}).otherwise({
		redirectTo: '/subscriptions'
	});
}]);
