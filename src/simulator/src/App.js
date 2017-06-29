import React, { Component } from 'react';
import './App.css';

import HeartRateSim from './components/heartrate';
import PedometerSim from './components/Pedometer';
import BloodPressureSim from './components/BloodPressure';

class App extends Component {
  render() {
    return (
      <div className="App">
        <div className="container-fluid">
          <br/>
          <br/>
          <h1>Simulation Panel</h1>
          <div className="row">
            <div className="col-md-6 center comp">
              <HeartRateSim />
            </div>
            <div className="col-md-6 center comp">
              <PedometerSim />
            </div>
          </div>
          <br/>
          <br/>
          <div className="row">
            <div className="col-md-6 center comp">
              <BloodPressureSim />
            </div>
          </div>
        </div>
      </div>
    );
  }
}

export default App;
