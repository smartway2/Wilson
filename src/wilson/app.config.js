(function() {
  'use strict';

  angular.module('Wilson').config(config)

  config.$inject = ['$stateProvider', '$urlRouterProvider', '$locationProvider'];

  angular.module('Wilson').run(function($rootScope) {
    $rootScope.$on('$stateChangeError', function(event, toState, toParams, fromState, fromParams, error){
      console.error(error);
    });
  });

  function config($stateProvider, $urlRouterProvider, $locationProvider) {
    $locationProvider.html5Mode(true)
    $stateProvider
      .state({
        name: 'home',
        url: '/',
        component: 'wilson'
      })

  }

}());
