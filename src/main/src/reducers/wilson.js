let initialState = {

}

export default (state = initialState, action) => {
    switch (action.type) {
        case 'INCREASE_ACTIVITY':
            return state;
        default:
            return state;
    }
};
