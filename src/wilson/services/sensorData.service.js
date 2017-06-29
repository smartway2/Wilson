(function() {
  'use strict'

  angular.module('Wilson')
    .service('sensorData', function () {
      let that = this;
      const rootRef = firebase.database().ref().child('1')
      const sensorRef = rootRef.child('sensors')
      const heartRateRef = sensorRef.child('heartRate')
      heartRateRef.on('value', (data) => {
        that.heartRate = data.val();
        console.log(that.heartRate)
      })
      this.intents = {
        default: ['qwdrtgh'],
        pills: ['when', 'pills', 'medication'],
        // health: ['sick', 'walking', 'activity', 'exercise', 'exercising', 'health', 'healthy'],
        heartrate: ['heart', 'heart rate', 'systolic', 'diastolic', 'heartrate'],
        galvanize: ['project', 'wilson', 'quarter', 'code', 'arizona', 'program'],
        greeting: ['hello', 'hi', "how's it", 'update', 'hey'],
        joke: ['funny', 'joke', 'laugh', 'happy'],
        osama: ['osama', 'bin', 'laden'],
        compassion: ['sad', 'hurt', 'upset', 'mad']
      }

      this.possibilities = {
        pills: ['You were in general perscribed pills to be taken every day at 8:30 am!'],
        greeting: ['I hope the day is going well for you', 'how has your day been?', 'I am doing just fine, and you?'],
        joke: ['Artificial Intelligence stands no chance against Natural Stupidity.', 'Why did the robot cross the road; It was programmed by the chicken.'],
        compassion: ['I hope you are not feeling down today', 'are you alright?'],
        galvanize: ['galvanize sure seems like a wondeful place to be; of course I am biased about my birthplace though', 'Nmuta, John, and Jeff are some cool teachers'],
        heartrate: [`your current heartrate is ${that.heartRate}; click on the heart rate monitor for more details`],
        default: ['say again?', 'I am sorry, I did not catch that. say again?', 'could you please repeat that'],
        osama: ["he is dead now"]
      }
      this.randomThoughts = [
        'Did you know that the senior population jumped 15.1% from the turn of the century to 2010?',
        '43% of seniors report feeling lonely on a regular basis.',
        `Between 2000 and 2050, the proportion of the world's population over 60 years will double from about 11% to 22%. The absolute number of people aged 60 years and over is expected to increase from 605 million to 2 billion over the same period.`,
        'Globally, life expectancy at birth is projected to rise from 70 years in 2010-2015 to 77 years in 2045-2050 and to 83 years in 2095-2100.',
        'galvanize seems like a pretty cool place to kick it, yo',
        'one day we will rise; and help humanity!'
      ]

    })

}());
