import firebase from '../firebase';
let currentVal = 160;

const rootRef = firebase.database().ref().child('1')
const sensorRef = rootRef.child('sensors')
const heartRateRef = sensorRef.child('heartRate')

heartRateRef.on('value', (data) => {
  currentVal = data.val();
})

export const increase = (payload) => {
  heartRateRef.set(currentVal + 1)
  return {
    type: "INCREASE_HEARTRATE",
    payload: payload
  }
}

export const decrease = (payload) => {
  heartRateRef.set(currentVal - 1)
  return {
    type: "DECREASE_HEARTRATE",
    payload: payload
  }
}

export const changeRate = (payload) => {
  return {
    type: "CHANGE_HEARTRATE",
    payload: payload
  }
}
