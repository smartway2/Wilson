import React, { Component } from 'react';
import '../App.css';
// import firebase from '../firebase';

import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import * as wilsonActions from '../actions/wilsonActions';
import Speech from 'react-speech';

class Wilson extends Component {

  render() {
    return (
      <div>
        <div>
          <h1>Wilson</h1>
          <Speech text="Welcome to react speech" />
        </div>
      </div>
    );
  }
}

function mapStateToProps(state) {
  console.log(state)
    return {
        state: state
    };
}

function mapDispatchToProps(dispatch) {
    return {
        actions: bindActionCreators(wilsonActions, dispatch)
    }
}

export default connect(mapStateToProps, mapDispatchToProps)(Wilson);
