(function() {
  'use strict'

  const app = angular.module('Wilson');

  app.component('diastolic',{
    bindings: {
      diastolic: '<'
    },
    templateUrl: 'components/diastolic.html',
    controller: function($http, $scope){
      const vm = this;

      vm.$onInit = function(){

      }
    },
    controllerAs: 'vm'
  });

}());
