import React, { Component } from 'react';
import '../App.css';
import firebase from '../firebase';

import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import * as pedometerActions from '../actions/pedometerActions';

class PedometerSim extends Component {

  componentDidMount(){
    const rootRef = firebase.database().ref().child('1')
    const sensorRef = rootRef.child('sensors')
    const pedometerRef = sensorRef.child('pedometer')
    pedometerRef.on('value', (data) => {
      this.props.actions.changeRate(data.val())
    })
  }

  render() {
    return (
      <div>
        <div>
          <h1>Pedometer</h1>
        </div>
        <div>
          <h3>{this.props.steps}</h3> steps/hour
        </div>
        <br />
        <script>console.log('window')</script>
        <button onClick={this.props.actions.increase}>INCREASE ACTIVITY</button>
        <button onClick={this.props.actions.decrease}>DECREASE ACTIVITY</button>
      </div>
    );
  }
}

function mapStateToProps(state) {
  console.log(state)
    return {
        steps: state.pedometer.steps
    };
}

function mapDispatchToProps(dispatch) {
    return {
        actions: bindActionCreators(pedometerActions, dispatch)
    }
}

export default connect(mapStateToProps, mapDispatchToProps)(PedometerSim);
