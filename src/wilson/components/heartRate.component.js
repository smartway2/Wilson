(function() {
  'use strict'

  const app = angular.module('Wilson');

  app.component('heartRate',{
    bindings: {
      heartRate: '<'
    },
    templateUrl: 'components/heartRate.html',
    controller: function($http, $scope){
      const vm = this;

      vm.$onInit = function(){
        
      }
    },
    controllerAs: 'vm'
  });

}());
