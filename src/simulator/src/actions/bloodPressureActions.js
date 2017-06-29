import firebase from '../firebase';
let currentSystolic = 117;
let currentDiastolic = 70;

const rootRef = firebase.database().ref().child('1');
const sensorRef = rootRef.child('sensors');
const bloodPressureRef = sensorRef.child('bloodPressure');
const systolicRef = bloodPressureRef.child('systolic');
const diastolicRef = bloodPressureRef.child('diastolic');

systolicRef.on('value', (data) => {
  currentSystolic = data.val();
});

diastolicRef.on('value', (data) => {
  currentDiastolic = data.val();
});

export const increaseSystolic = (payload) => {
  systolicRef.set(currentSystolic + 1)
  return {
    type: "INCREASE_SYSTOLIC",
    payload: payload
  }
}

export const decreaseSystolic = (payload) => {
  systolicRef.set(currentSystolic - 1)
  return {
    type: "DECREASE_SYSTOLIC",
    payload: payload
  }
}

export const changeSystolic = (payload) => {
  return {
    type: "CHANGE_SYSTOLIC",
    payload: payload
  }
}

export const increaseDiastolic = (payload) => {
  diastolicRef.set(currentDiastolic + 1)
  return {
    type: "INCREASE_DIASTOLIC",
    payload: payload
  }
}

export const decreaseDiastolic = (payload) => {
  diastolicRef.set(currentDiastolic - 1)
  return {
    type: "DECREASE_DIASTOLIC",
    payload: payload
  }
}

export const changeDiastolic = (payload) => {
  return {
    type: "CHANGE_DIASTOLIC",
    payload: payload
  }
}
