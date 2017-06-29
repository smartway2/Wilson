let initialState = {
  steps: 466
}

export default (state = initialState, action) => {
    switch (action.type) {
        case 'INCREASE_ACTIVITY':
            return state;
        case 'DECREASE_ACTIVITY':
            return state;
        case 'CHANGE_ACTIVITY':
            return {
              ...state,
              steps: action.payload
            }
        default:
            return state;
    }
};
