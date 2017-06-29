(function() {
  'use strict'

  const app = angular.module('Wilson');

  app.component('wilson',{
    templateUrl: 'components/wilson.html',
    controller: function($http, sensorData, $scope){

      const vm = this;



      vm.$onInit = function(){
        vm.finalTranscript = '';
        vm.initializeWebKit()
        const rootRef = firebase.database().ref().child('1');
        const sensorRef = rootRef.child('sensors');
        const heartRateRef = sensorRef.child('heartRate');
        const pedometerRef = sensorRef.child('pedometer');
        const bpRef = sensorRef.child('bloodPressure');
        const sysRef = bpRef.child('systolic');
        const diaRef = bpRef.child('diastolic');
        heartRateRef.on('value', (data) => {
          vm.heartRate = data.val();
          $scope.$apply();
          vm.checkVitals();
        })
        pedometerRef.on('value', (data) => {
          vm.pedometer = data.val();
          $scope.$apply();
          vm.checkVitals();
        })
        sysRef.on('value', (data) => {
          vm.systolic = data.val();
          $scope.$apply();
          vm.checkVitals();
        })
        diaRef.on('value', (data) => {
          vm.diastolic = data.val();
          $scope.$apply();
          vm.checkVitals();
        })
        setInterval(vm.thinkOutLoud(), 50000)

        function randomBetween(range) {
          	var min = range[0],
                max = range[1];
            if (min < 0) {
                return min + Math.random() * (Math.abs(min)+max);
            }else {
                return min + Math.random() * max;
            }
        }

        $.fn.equalizerAnimation = function(speed, barsHeight){
            var $equalizer = $(this);
            setInterval(function(){
              	$equalizer.find('span').each(function(i){
                  $(this).css({ height:randomBetween(barsHeight[i])+'px' });
                });
            },speed);
            $equalizer.toggleClass('paused');
        }

        var barsHeight = [
          [40, 80],
          [60, 140],
          [80, 160],
          [60, 160],
          [60, 120]
        ];
        $('.equalizer').equalizerAnimation(180, barsHeight);

      }

      vm.meds = function(){
        $('.modal-content').empty()
        $.getJSON('http://localhost:8080/api/v1/medication/1', (data) => {
          for(let x of data.rows){
            console.log(data)
            $('.modal-content').append(`<br><li>${JSON.stringify(x)}</li>`)
          }
          $('#heartRateModal').modal('show');
        })
      }

      vm.showModal = function(modal, chart){
        $('#' + modal ).modal('show');
        vm.renderChart(chart);
      }
      vm.alerted = {
        heartRateLow: false,
        heartRateHigh: false,
        sysBPhigh: false,
        sysBPlow: false,
        diaBPhigh: false,
        diaBPlow: false
      }
      vm.checkVitals = function(){
        if(vm.heartRate > 130 && vm.pedometer < 660){
          let speech = `At ${vm.heartRate}, your heart rate is unusually high for your current level of activity and cardiac history. I have contacted all of your emergency contacts for immediate assistance.`
          if(!vm.alerted.heartRateHigh){
            vm.executeSpeech(speech);
          }
          vm.alerted.heartRateHigh = true;
        } else if (vm.heartRate < 50){
          let speech = `At ${vm.heartRate}, your heart rate is unusually low, given your cardiac history. I have contacted some emergency contacts to check on you`
          if(!vm.alerted.heartRatelow){
            vm.executeSpeech(speech);
          }
          vm.alerted.heartRatelow = true;
        } else if (vm.systolic > 160){
          let speech = `your systolic blood pressure is now dangerously high at ${vm.systolic}. appropriate emergency contacts have been notified.`
          if(!vm.alerted.sysBPhigh){
            vm.executeSpeech(speech);
          }
          vm.alerted.sysBPhigh = true;
        } else if (vm.diastolic > 100){
          let speech = `at ${vm.diastolic}, your diastolic blood pressure is symptomatic of stage 2 hypertension. emergency contacts have been notified.`
          if(!vm.alerted.diaBPhigh){
            vm.executeSpeech(speech);
          }
          vm.alerted.diaBPhigh = true;
        } else if (vm.diastolic < 55 || vm.systolic < 80){
          let speech = `your blood pressure is getting unusually low. I have asked your contacts to check up on you.`
          if(!vm.alerted.diaBPlow){
            vm.executeSpeech(speech);
          }
          vm.alerted.diaBPlow = true;
        }
      }

      vm.thinkOutLoud = function(){
        vm.executeSpeech(sensorData.randomThoughts[Math.floor(Math.random() * sensorData.randomThoughts.length)])
      }

      vm.renderChart = function(chart){
        let chartData = {
          heartRateChart: {
            data: [60, 72, vm.heartRate],
            labels: ['Overall Average Heart Rate', 'Your Average', 'Current'],
            title: 'Heart Rate (beats/minute)',
            type: 'horizontalBar',
            speech: function(){
              if(vm.heartRate < 60){
                return `your heart rate is ${vm.heartRate} beats per minute, which is well below the average for an adult your age.`
              } else if (vm.heartRate > 59 && vm.heartRate < 100){
                return `at ${vm.heartRate} beats per minute, your heart rate falls within the average range for an adult your age, keep it up!`
              } else if (vm.heartRate > 99 && vm.pedometer > 660){
                return `you have been fairly active in the last few hours, so although your heart rate is high, at ${vm.heartRate} beats per minute, it is normal.`
              } else if (vm.heartRate > 99 && vm.heartRate < 110){
                return `your heart current rate is above the average for your age range, and you may want to consider seeing a doctor soon`
              } else {
                return `at ${vm.heartRate} beats per minute, your heart rate is high given your current level of activity. I recommend consulting a health expert`
              }
            }
          },
          pedometer: {
            data: [7500, 6990, vm.pedometer * 15],
            labels: ['Overall Average Steps', 'Your Average Steps', 'Current Pace'],
            title: 'Steps (steps/day)',
            type: 'bar',
            speech: function(){
              if(vm.pedometer < 430){
                return `you are currently averaging ${vm.pedometer} steps per hour, which means you have been more sedintary than usual, and you are not on pace to meet your average`
              } else if (vm.pedometer > 430 && vm.pedometer < 500){
                return `you are averaging ${vm.pedometer} steps per hour, and on pace for ${vm.pedometer * 15} steps today; which is about average for you`;
              } else {
                return `you are averaging ${vm.pedometer} steps per hour, and you are on pace to exceed your goal for the day!`
              }
            }
          },
          systolic: {
            data: [120, 117, vm.systolic],
            labels: ['Overall Average Systolic Blood Pressure', 'Your Average', 'Current'],
            title: 'Systolic Blood Pressure (mm Hg)',
            type: 'horizontalBar',
            speech: function(){
              if(vm.systolic < 121 && vm.systolic > 90){
                return `at ${vm.systolic} your systolic blood pressure falls within the average range for you and your health history`
              } else if (vm.systolic < 91){
                return `your systolic blood pressure is low. it is ${vm.systolic}, and you may want to consult a health professional`
              } else if (vm.systolic > 120 && vm.systolic < 140){
                return `your systolic blood pressure is a little bit higher than normal. it could simply be something you ate this morning. make sure to make healthy food choices!`
              } else if (vm.systolic > 139 && vm.systolic < 160){
                return `at ${vm.systolic}, your systolic blood pressure is symptomatic of stage 1 hypertension. consider consulting a physician`;
              } else {
                return `your systolic blood pressure, at ${vm.systolic}, is well above normal and in stage 2 hypertension. consult a health professional`
              }
            }
          },
          diastolic: {
            data: [80, 74, vm.diastolic],
            labels: ['Overall Average diastolic Blood Pressure', 'Your Average', 'Current'],
            title: 'Diastolic Blood Pressure(mm Hg)',
            type: 'bar',
            speech: function(){
              if(vm.diastolic < 81 && vm.diastolic > 59){
                return `at ${vm.diastolic} your diastolic blood pressure falls within the average range for you and your health history`
              } else if (vm.diastolic < 60){
                return `your diastolic blood pressure is low. it is ${vm.diastolic}, and you may want to consult a health professional`
              } else if (vm.diastolic > 80 && vm.diastolic < 90){
                return `your diastolic blood pressure is a little bit higher than normal. it could simply be something you ate this morning. make sure to make healthy food choices!`
              } else if (vm.diastolic > 89 && vm.diastolic < 100){
                return `at ${vm.diastolic}, your diastolic blood pressure is symptomatic of stage 1 hypertension. consider consulting a physician`;
              } else {
                return `your diastolic blood pressure, at ${vm.diastolic}, is well above normal and in stage 2 hypertension. consult a health professional`
              }
            }
          },
        }
        $('.modal-content').empty()
        $('.modal-content').append(`<canvas id="heartRateChart" width="150" height="150"></canvas>`)
        var ctx = document.getElementById('heartRateChart').getContext('2d');
        var myChart = new Chart(ctx, {
            type: chartData[chart].type,
            data: {
                labels: chartData[chart].labels,
                datasets: [{
                    label: chartData[chart].title,
                    data: chartData[chart].data,
                    backgroundColor: [
                        'rgba(255, 99, 132, 0.2)',
                        'rgba(54, 162, 235, 0.2)',
                        'rgba(255, 206, 86, 0.2)',
                        'rgba(75, 192, 192, 0.2)',
                        'rgba(153, 102, 255, 0.2)',
                        'rgba(255, 159, 64, 0.2)'
                    ],
                    borderColor: [
                        'rgba(255,99,132,1)',
                        'rgba(54, 162, 235, 1)',
                        'rgba(255, 206, 86, 1)',
                        'rgba(75, 192, 192, 1)',
                        'rgba(153, 102, 255, 1)',
                        'rgba(255, 159, 64, 1)'
                    ],
                    borderWidth: 1
                }]
            },
            options: {
                scales: {
                    yAxes: [{
                        ticks: {
                            beginAtZero:true
                        }
                    }],
                    xAxes: [{
                        ticks: {
                            beginAtZero:true
                        }
                    }]
                }
            }
        });
        vm.executeSpeech(chartData[chart].speech())
      }

      vm.executeSpeech = function(str){
        responsiveVoice.speak("" + str);
        $('.equalizer').toggleClass('paused', null, true);
        setTimeout(vm.resolveSpeech, 2000);
      }

      vm.resolveSpeech = function(){
        if(responsiveVoice.isPlaying()){
          setTimeout(vm.resolveSpeech, 2500)
        } else {
          $('.equalizer').toggleClass('paused', null, false);
          return;
        }
      }

      vm.executeReception = function(carrier){
        console.log('fired!!', carrier)
        vm.finalTranscript = carrier;
        $scope.$apply()
      }

      vm.executeReceptionFinal = function(carrier, interim){
        console.log('fired!!', carrier)
        vm.finalTranscript = carrier;
        vm.interpretText(carrier);
        $scope.$apply()
      }

      vm.interpretText = function(text){
        let scores = {};
        let maxScore = 'default';

        for (let x in sensorData.intents){
          sensorData.intents[x].forEach(y => {
            if(!(x in scores)){
              scores[x] = 0;
            }
            const reg = new RegExp (y, 'ig');
            reg.test(text) ? scores[x] += 1 : null;
            if(scores[x] > scores[maxScore]){
              maxScore = x;
            }
          })
        }
        console.log(scores);
        var response = sensorData.possibilities[maxScore][Math.floor(Math.random() * sensorData.possibilities[maxScore].length)];
        vm.executeSpeech(response);

      }

      vm.initializeWebKit = function(){
        console.log(sensorData.houses)
          let final_transcript = '';
          let create_email = false;
          let recognizing = false;
          let ignore_onend;
          let start_timestamp;
          if (!('webkitSpeechRecognition' in window)) {
            upgrade();
          } else {
            start_button.style.display = 'inline-block';
            var recognition = new webkitSpeechRecognition();
            recognition.continuous = true;
            recognition.interimResults = true;

            recognition.onstart = function() {
              recognizing = true;
            };

            recognition.onerror = function(event) {
              if (event.error == 'no-speech') {
                ignore_onend = true;
              }
              if (event.error == 'audio-capture') {
                ignore_onend = true;
              }
              if (event.error == 'not-allowed') {
                if (event.timeStamp - start_timestamp < 100) {
                } else {
                }
                ignore_onend = true;
              }
            };

            recognition.onend = function() {
              recognizing = false;
              if (ignore_onend) {
                return;
              }
              if (!final_transcript) {
                return;
              }
            };

            recognition.onresult = function(event) {
              var interim_transcript = '';
              console.log(event.results)
              for (let i = event.resultIndex; i < event.results.length; ++i) {
                if (event.results[i].isFinal) {
                  final_transcript += event.results[i][0].transcript;
                  console.log('final', final_transcript);
                  vm.executeReceptionFinal(final_transcript, interim_transcript)
                } else {
                  interim_transcript += event.results[i][0].transcript;
                  vm.executeReception(interim_transcript)
                }
              }
              final_transcript = capitalize(final_transcript);
              console.log('final', final_transcript);
            };
          }

          let current_style;
          function showButtons(style) {
            if (style == current_style) {
              return;
            }
            current_style = style;
          }

          function startButton(event) {
            if (recognizing) {
              recognition.stop();
              return;
            }
            final_transcript = '';
            recognition.lang = 'en-US';
            recognition.start();
            ignore_onend = false;
            showButtons('none');
            start_timestamp = event.timeStamp;
          }
          function upgrade() {
            start_button.style.visibility = 'hidden';
          }
          let two_line = /\n\n/g;
          let one_line = /\n/g;
          function linebreak(s) {
            return s.replace(two_line, '<p></p>').replace(one_line, '<br>');
          }

          let first_char = /\S/;
          function capitalize(s) {
            return s.replace(first_char, function(m) { return m.toUpperCase(); });
          }
          $('#start_button').click(startButton)
      }
    },
    controllerAs: 'vm'
  });



}());
