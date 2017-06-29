let initialState = {
  systolic: 117,
  diastolic: 70
}

export default (state = initialState, action) => {
    switch (action.type) {
        case 'INCREASE_SYSTOLIC':
            return state;
        case 'DECREASE_SYSTOLIC':
            return state;
        case 'CHANGE_SYSTOLIC':
            return {
              ...state,
              systolic: action.payload
            }
        case 'INCREASE_DIASTOLIC':
            return state;
        case 'DECREASE_DIASTOLIC':
            return state;
        case 'CHANGE_DIASTOLIC':
            return {
              ...state,
              diastolic: action.payload
            }
        default:
            return state;
    }
};
