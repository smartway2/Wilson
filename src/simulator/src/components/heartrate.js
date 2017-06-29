import React, { Component } from 'react';
import '../App.css';
import firebase from '../firebase';

import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import * as heartRateActions from '../actions/heartRateActions';

class HeartRateSim extends Component {

  componentDidMount(){
    const rootRef = firebase.database().ref().child('1')
    const sensorRef = rootRef.child('sensors')
    const heartRateRef = sensorRef.child('heartRate')
    heartRateRef.on('value', (data) => {
      this.props.actions.changeRate(data.val())
    })
  }

  render() {
    return (
      <div>
        <div>
          <h1>Heart Rate Monitor</h1>
        </div>
        <div>
          <h3>{this.props.heartRate.heartRate}</h3> heartbeats per minute
        </div>
        <br />
        <button onClick={this.props.actions.increase}>INCREASE HEARTRATE</button>
        <button onClick={this.props.actions.decrease}>DECREASE HEARTRATE</button>
      </div>
    );
  }
}

function mapStateToProps(state) {
    return {
        heartRate: state.heartRate
    };
}

function mapDispatchToProps(dispatch) {
    return {
        actions: bindActionCreators(heartRateActions, dispatch)
    }
}

export default connect(mapStateToProps, mapDispatchToProps)(HeartRateSim);
