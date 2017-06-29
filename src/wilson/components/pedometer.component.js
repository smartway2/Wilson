(function() {
  'use strict'

  const app = angular.module('Wilson');

  app.component('pedometer',{
    bindings: {
      pedometer: '<'
    },
    templateUrl: 'components/pedometer.html',
    controller: function($http, $scope){
      const vm = this;

      vm.$onInit = function(){

      }
    },
    controllerAs: 'vm'
  });

}());
