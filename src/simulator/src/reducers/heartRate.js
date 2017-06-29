let initialState = {
  heartRate: 160
}

export default (state = initialState, action) => {
    switch (action.type) {
        case 'INCREASE_HEARTRATE':
            return state;
        case 'DECREASE_HEARTRATE':
            return state;
        case 'CHANGE_HEARTRATE':
            return {
              ...state,
              heartRate: action.payload
            }
        default:
            return state;
    }
};
