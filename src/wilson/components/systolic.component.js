(function() {
  'use strict'

  const app = angular.module('Wilson');

  app.component('systolic',{
    bindings: {
      systolic: '<'
    },
    templateUrl: 'components/systolic.html',
    controller: function($http, $scope){
      const vm = this;

      vm.$onInit = function(){

      }
    },
    controllerAs: 'vm'
  });

}());
