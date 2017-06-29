import firebase from '../firebase';
let currentVal = 466;

const rootRef = firebase.database().ref().child('1')
const sensorRef = rootRef.child('sensors')
const pedometerRef = sensorRef.child('pedometer')

pedometerRef.on('value', (data) => {
  currentVal = data.val();
})

export const increase = (payload) => {
  pedometerRef.set(currentVal + 2)
  return {
    type: "INCREASE_ACTIVITY",
    payload: payload
  }
}

export const decrease = (payload) => {
  pedometerRef.set(currentVal - 2)
  return {
    type: "DECREASE_ACTIVITY",
    payload: payload
  }
}

export const changeRate = (payload) => {
  return {
    type: "CHANGE_ACTIVITY",
    payload: payload
  }
}
