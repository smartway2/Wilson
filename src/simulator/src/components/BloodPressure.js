import React, { Component } from 'react';
import '../App.css';
import firebase from '../firebase';

import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import * as bloodPressureActions from '../actions/bloodPressureActions';

class BloodPressureSim extends Component {

  componentDidMount(){
    const rootRef = firebase.database().ref().child('1');
    const sensorRef = rootRef.child('sensors');
    const bloodPressureRef = sensorRef.child('bloodPressure');
    const systolicRef = bloodPressureRef.child('systolic');
    const diastolicRef = bloodPressureRef.child('diastolic');

    systolicRef.on('value', (data) => {
      this.props.actions.changeSystolic(data.val())
    });

    diastolicRef.on('value', (data) => {
      this.props.actions.changeDiastolic(data.val())
    });
  }

  render() {
    return (
      <div>
        <div>
          <h1>BloodPressure</h1>
        </div>
        <div>
          <h3>{this.props.systolic}</h3> Systolic (mm Hg)
        </div>
        <br />
        <button onClick={this.props.actions.increaseSystolic}>INCREASE SYSTOLIC</button>
        <button onClick={this.props.actions.decreaseSystolic}>DECREASE SYSTOLIC</button>
        <div>
          <h3>{this.props.diastolic}</h3> Diastolic(mm Hg)
        </div>
        <br />
        <button onClick={this.props.actions.increaseDiastolic}>INCREASE DIASTOLIC</button>
        <button onClick={this.props.actions.decreaseDiastolic}>DECREASE DIASTOLIC</button>
      </div>
    );
  }
}

function mapStateToProps(state) {
  console.log(state)
    return {
        systolic: state.bloodPressure.systolic,
        diastolic: state.bloodPressure.diastolic
    };
}

function mapDispatchToProps(dispatch) {
    return {
        actions: bindActionCreators(bloodPressureActions, dispatch)
    }
}

export default connect(mapStateToProps, mapDispatchToProps)(BloodPressureSim);
