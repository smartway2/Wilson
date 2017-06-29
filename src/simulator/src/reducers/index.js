import heartRate from './heartRate';
import pedometer from './pedometer'
import bloodPressure from './bloodPressure'
import { combineReducers } from 'redux';

const rootReducer = combineReducers({
    heartRate,
    pedometer,
    bloodPressure
});
export default rootReducer;
